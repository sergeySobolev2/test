package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	// SERVICE_TOKEN - токен для авторизации асинхронного сервиса
	SERVICE_TOKEN = "async-partition-soundproof-key"
)

// AsyncResultRequest - структура для приема результата от асинхронного сервиса
type AsyncResultRequest struct {
	RequestID        uint    `json:"calculation_id" binding:"required"`
	NoiseReductionDB float64 `json:"noise_reduction_db" binding:"required"`
	RoomArea         float64 `json:"room_area" binding:"required"`
	ServiceToken     string  `json:"service_token" binding:"required"`
}

// TriggerCalculationRequest - структура для запроса на расчет
type TriggerCalculationRequest struct {
	RoomArea    float64 `json:"room_area" binding:"required,gt=0"`
	ThicknessCM float64 `json:"thickness_cm" binding:"required,gt=0"`
}

// ApiReceiveAsyncResult - прием результата расчета от асинхронного сервиса
// @Summary Принять результат асинхронного расчета
// @Description Метод для приема результата от асинхронного сервиса. Требует токен авторизации сервиса.
// @Tags async
// @Accept json
// @Produce json
// @Param request body AsyncResultRequest true "Результат расчета"
// @Success 200 {object} object{status=string,calculation_id=uint}
// @Failure 400 {object} object{error=string}
// @Failure 401 {object} object{error=string}
// @Failure 404 {object} object{error=string}
// @Router /api/async/result [post]
func (h *Handler) ApiReceiveAsyncResult(ctx *gin.Context) {
	var req AsyncResultRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Проверка токена авторизации сервиса
	if req.ServiceToken != SERVICE_TOKEN {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid service token"})
		return
	}

	// Получаем заявку для проверки существования
	calc, err := h.Repository.GetCalculationByID(req.RequestID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "calculation not found"})
		return
	}

	// Обновляем снижение шума и инкрементируем счетчик расчетов
	fields := map[string]interface{}{
		"noise_reduction_db": req.NoiseReductionDB,
		"calculations_count": calc.CalculationsCount + 1,
	}

	if err := h.Repository.UpdateCalculationFields(req.RequestID, fields); err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":         "success",
		"calculation_id": req.RequestID,
		"message":        "Noise reduction updated successfully",
	})
}

// ApiTriggerAsyncCalculation - запуск асинхронного расчета
// @Summary Запустить асинхронный расчет
// @Description Отправляет запрос на расчет снижения шума в асинхронный сервис с указанием площади помещения
// @Tags async
// @Security BearerAuth
// @Security CookieAuth
// @Accept json
// @Produce json
// @Param id path int true "ID заявки"
// @Param request body TriggerCalculationRequest true "Площадь помещения"
// @Success 202 {object} object{status=string,message=string}
// @Failure 400 {object} object{error=string}
// @Failure 401 {object} object{error=string}
// @Failure 404 {object} object{error=string}
// @Router /api/requests/{id}/calculate [post]
func (h *Handler) ApiTriggerAsyncCalculation(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	// Получаем площадь из запроса
	var req TriggerCalculationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Получаем заявку
	request, err := h.Repository.GetCalculationByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "calculation not found"})
		return
	}

	// Сохраняем площадь и толщину в БД
	fields := map[string]interface{}{
		"room_area":          req.RoomArea,
		"required_thickness": req.ThicknessCM,
	}
	if err := h.Repository.UpdateCalculationFields(uint(id), fields); err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	// Получаем материалы перегородок в заявке
	materials := make([]string, 0)
	if details, err := h.Repository.GetCalculationWithPartitions(uint(id)); err == nil {
		for _, p := range details.Partitions {
			if p.Material != "" {
				materials = append(materials, p.Material)
			}
		}
	}

	// Отправляем запрос в асинхронный сервис
	asyncServiceURL := os.Getenv("ASYNC_SERVICE_URL")
	if asyncServiceURL == "" {
		asyncServiceURL = "http://localhost:8083"
	}

	payload := map[string]interface{}{
		"calculation_id": request.ID,
		"room_area":      req.RoomArea,
		"thickness_cm":   req.ThicknessCM,
		"materials":      materials,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	url := fmt.Sprintf("%s/calculate", asyncServiceURL)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "async service error"})
		return
	}

	ctx.JSON(http.StatusAccepted, gin.H{
		"status":  "accepted",
		"message": "Calculation request sent to async service",
	})
}

// triggerAsyncCalculationInternal - внутренняя функция для запуска асинхронного расчета
// Используется при автоматическом запуске расчета после формирования заявки
func (h *Handler) triggerAsyncCalculationInternal(requestID uint, RoomArea float64) {
	asyncServiceURL := os.Getenv("ASYNC_SERVICE_URL")
	if asyncServiceURL == "" {
		asyncServiceURL = "http://localhost:8083"
	}

	payload := map[string]interface{}{
		"calculation_id": requestID,
		"room_area":      RoomArea,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return
	}

	url := fmt.Sprintf("%s/calculate", asyncServiceURL)
	_, _ = http.Post(url, "application/json", bytes.NewBuffer(jsonData))
}

package handler

import (
	"net/http"
	"strconv"
	"time"

	"partitionlab/internal/app/ds"
	"partitionlab/internal/app/middleware"

	"github.com/gin-gonic/gin"
)

// ApiGetCart получить корзину (черновик заявки)
// @Summary Получить корзину
// @Description Требуется авторизация. Возвращает ID черновика и количество симптомов в нем.
// @Tags requests
// @Security BearerAuth
// @Security CookieAuth
// @Accept json
// @Produce json
// @Success 200 {object} object{data=object{draft_id=uint,count=int},total=int,filters=object}
// @Failure 401 {object} object{error=string}
// @Router /api/requests/cart [get]
func (h *Handler) ApiGetCart(ctx *gin.Context) {
	userID, ok := middleware.GetCurrentUserID(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	count, _ := h.Repository.CountDraftItems(userID)
	draft, _ := h.Repository.GetDraftCalculation(userID)
	var draftID *uint
	if draft != nil {
		draftID = &draft.ID
	}
	jsonResponse(ctx, gin.H{"draft_id": draftID, "count": count}, 1, gin.H{})
}

// ApiListRequests получить список заявок
// @Summary Получить список заявок
// @Description Требуется авторизация. Обычный пользователь видит только свои заявки. Модератор видит все заявки.
// @Tags requests
// @Security BearerAuth
// @Security CookieAuth
// @Accept json
// @Produce json
// @Param status query string false "Фильтр по статусу"
// @Param from query string false "Дата начала (formed_at)"
// @Param to query string false "Дата окончания (formed_at)"
// @Success 200 {object} object{data=[]object,total=int,filters=object}
// @Failure 401 {object} object{error=string}
// @Router /api/requests [get]
func (h *Handler) ApiListCalculations(ctx *gin.Context) {
	userID, ok := middleware.GetCurrentUserID(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Проверяем, является ли пользователь модератором
	isModerator := middleware.IsCurrentUserModerator(ctx)

	status := ctx.Query("status")
	from := ctx.Query("from")
	to := ctx.Query("to")
	var statusPtr *string
	if status != "" {
		statusPtr = &status
	}

	var list []ds.Calculation
	var err error

	if isModerator {
		// Модератор видит все заявки
		list, err = h.Repository.ListAllRequestsWithFilters(statusPtr, &from, &to)
	} else {
		// Обычный пользователь видит только свои заявки
		list, err = h.Repository.ListRequestsWithFilters(userID, statusPtr, &from, &to)
	}

	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	// Добавляем вычисляемое поле: количество симптомов с непустым комментарием в заявке
	type flatRequestItem struct {
		ID                uint       `json:"id"`
		UserID            uint       `json:"user_id"`
		UserLogin         string     `json:"user_login"`
		Status            string     `json:"status"`
		FormedAt          *time.Time `json:"formed_at"`
		CompletedAt       *time.Time `json:"completed_at"`
		ModeratorID       *uint      `json:"moderator_id"`
		ModeratorLogin    *string    `json:"moderator_login"`
		RoomArea          *float64   `json:"room_area"`
		NoiseReductionDB  *float64   `json:"noise_reduction_db"`
		RequiredThickness *float64   `json:"required_thickness"`
		ExpertComment     string     `json:"expert_comment"`
		CommentsCount     int64      `json:"comments_count"`
		CalculationsCount int        `json:"calculations_count"` // Счетчик результатов расчета
	}
	resp := make([]flatRequestItem, 0, len(list))
	for _, r := range list {
		cnt, _ := h.Repository.CountCalculationPartitionsWithComment(r.ID)
		var modLogin *string
		if r.Moderator != nil {
			modLogin = &r.Moderator.Login
		}
		comment := ""
		if r.ExpertComment.Valid {
			comment = r.ExpertComment.String
		}
		resp = append(resp, flatRequestItem{
			ID:                r.ID,
			UserID:            r.UserID,
			UserLogin:         r.User.Login,
			Status:            r.Status,
			FormedAt:          r.FormedAt,
			CompletedAt:       r.CompletedAt,
			ModeratorID:       r.ModeratorID,
			ModeratorLogin:    modLogin,
			RoomArea:          r.RoomArea,
			NoiseReductionDB:  r.NoiseReductionDB,
			RequiredThickness: r.RequiredThickness,
			ExpertComment:     comment,
			CommentsCount:     cnt,
			CalculationsCount: r.CalculationsCount,
		})
	}
	jsonResponse(ctx, resp, int64(len(resp)), gin.H{"status": status, "from": from, "to": to})
}

// ApiGetRequest получить заявку по ID
// @Summary Получить заявку по ID
// @Description Требуется авторизация. Возвращает заявку с симптомами.
// @Tags requests
// @Security BearerAuth
// @Security CookieAuth
// @Accept json
// @Produce json
// @Param id path int true "ID заявки"
// @Success 200 {object} object{data=object,total=int,filters=object}
// @Failure 404 {object} object{error=string}
// @Failure 401 {object} object{error=string}
// @Router /api/requests/{id} [get]
func (h *Handler) ApiGetCalculation(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	dr, err := h.Repository.GetCalculationWithPartitions(uint(id))
	if err != nil || dr == nil {
		h.errorHandler(ctx, http.StatusNotFound, err)
		return
	}
	// скрываем удаленные
	if dr.Calculation.Status == "удален" || dr.Calculation.Status == "удалён" {
		ctx.Status(http.StatusNotFound)
		return
	}
	jsonResponse(ctx, dr, 1, gin.H{"id": id})
}

// ApiUpdateRequest обновить заявку
// @Summary Обновить заявку
// @Description Требуется авторизация. Доступно только владельцу черновика.
// @Tags requests
// @Security BearerAuth
// @Security CookieAuth
// @Accept json
// @Produce json
// @Param id path int true "ID заявки"
// @Param request body object{room_area=float64,expert_comment=string} true "Обновляемые поля"
// @Success 200 {object} object{data=ds.Calculation,total=int,filters=object}
// @Failure 400 {object} object{error=string}
// @Failure 401 {object} object{error=string}
// @Failure 403 {object} object{error=string}
// @Router /api/requests/{id} [put]
func (h *Handler) ApiUpdateCalculation(ctx *gin.Context) {
	userID, ok := middleware.GetCurrentUserID(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	idStr := ctx.Param("id")
	id64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	id := uint(id64)
	// только владелец
	if owner, err := h.Repository.IsCalculationOwner(userID, id); err != nil || !owner {
		h.errorHandler(ctx, http.StatusForbidden, err)
		return
	}
	rq, err := h.Repository.GetCalculationByID(id)
	if err != nil {
		h.errorHandler(ctx, http.StatusNotFound, err)
		return
	}
	if rq.Status != "черновик" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "only draft can be updated"})
		return
	}
	type bodyT struct {
		RoomArea      *float64 `json:"room_area"`
		ExpertComment *string  `json:"expert_comment"`
	}
	var body bodyT
	if err := ctx.ShouldBindJSON(&body); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	fields := map[string]interface{}{}
	if body.RoomArea != nil {
		fields["room_area"] = *body.RoomArea
	}
	if body.ExpertComment != nil {
		fields["expert_comment"] = *body.ExpertComment
	}
	if err := h.Repository.UpdateCalculationFields(id, fields); err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}
	rq, _ = h.Repository.GetCalculationByID(id)
	jsonResponse(ctx, rq, 1, gin.H{"id": id})
}

// ApiFormRequest сформировать заявку
// @Summary Сформировать заявку
// @Description Требуется авторизация. Переводит черновик в статус "сформирован". Доступно только владельцу.
// @Tags requests
// @Security BearerAuth
// @Security CookieAuth
// @Accept json
// @Produce json
// @Param id path int true "ID заявки"
// @Success 200 {object} object{data=ds.Calculation,total=int,filters=object}
// @Failure 400 {object} object{error=string}
// @Failure 401 {object} object{error=string}
// @Failure 403 {object} object{error=string}
// @Router /api/requests/{id}/form [put]
func (h *Handler) ApiFormCalculation(ctx *gin.Context) {
	userID, ok := middleware.GetCurrentUserID(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	idStr := ctx.Param("id")
	id64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	id := uint(id64)
	// только владелец может формировать черновик
	if owner, err := h.Repository.IsCalculationOwner(userID, id); err != nil || !owner {
		h.errorHandler(ctx, http.StatusForbidden, err)
		return
	}
	rq, err := h.Repository.GetCalculationByID(id)
	if err != nil {
		h.errorHandler(ctx, http.StatusNotFound, err)
		return
	}
	if rq.Status != "черновик" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "only draft can be formed"})
		return
	}
	// обязательные поля: есть хотя бы один симптом
	cnt, err := h.Repository.CountCalculationPartitions(id)
	if err != nil || cnt == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "empty request"})
		return
	}
	if err := h.Repository.SetFormed(id, time.Now()); err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	// Примечание: асинхронный расчет процента дегидратации будет выполнен с задержкой
	// Результат будет получен через endpoint /api/async/result
	if rq.RoomArea != nil && *rq.RoomArea > 0 {
		// Запускаем в горутине чтобы не блокировать ответ
		go h.triggerAsyncCalculationInternal(id, *rq.RoomArea)
	}

	rq, _ = h.Repository.GetCalculationByID(id)
	jsonResponse(ctx, rq, 1, gin.H{"id": id})
}

// ApiCompleteRequest завершить/отклонить заявку
// @Summary Завершить или отклонить заявку
// @Description Требуется авторизация модератора. Переводит заявку в статус "завершен" или "отклонен".
// @Tags requests
// @Security BearerAuth
// @Security CookieAuth
// @Accept json
// @Produce json
// @Param id path int true "ID заявки"
// @Param request body object{status=string,room_area=float64,noise_reduction_db=float64,expert_comment=string} true "Данные для завершения"
// @Success 200 {object} object{data=ds.Calculation,total=int,filters=object}
// @Failure 400 {object} object{error=string}
// @Failure 401 {object} object{error=string}
// @Failure 403 {object} object{error=string}
// @Router /api/requests/{id}/complete [put]
func (h *Handler) ApiCompleteCalculation(ctx *gin.Context) {
	userID, ok := middleware.GetCurrentUserID(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Проверка прав модератора (дополнительная проверка)
	if !middleware.IsCurrentUserModerator(ctx) {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "only moderator can complete/reject"})
		return
	}

	idStr := ctx.Param("id")
	id64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	id := uint(id64)
	type bodyT struct {
		Status            string   `json:"status" binding:"required"` // завершен | отклонен
		RoomArea          *float64 `json:"room_area"`
		NoiseReductionDB  *float64 `json:"noise_reduction_db"`
		RequiredThickness *float64 `json:"required_thickness"`
		ExpertComment     *string  `json:"expert_comment"`
	}
	var body bodyT
	if err := ctx.ShouldBindJSON(&body); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	if body.Status != "завершен" && body.Status != "отклонен" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid status"})
		return
	}
	rq, err := h.Repository.GetCalculationByID(id)
	if err != nil {
		h.errorHandler(ctx, http.StatusNotFound, err)
		return
	}
	if rq.Status != "сформирован" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "only formed can be completed"})
		return
	}
	var RequiredThickness *float64
	if body.Status == "завершен" {
		if body.RoomArea == nil || body.NoiseReductionDB == nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "weight and percent required"})
			return
		}
		if body.RequiredThickness != nil {
			RequiredThickness = body.RequiredThickness
		} else {
			RequiredThickness = rq.RequiredThickness
		}
	}
	if err := h.Repository.SetCompleted(id, userID, body.Status, time.Now(), body.RoomArea, body.NoiseReductionDB, RequiredThickness, body.ExpertComment); err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}
	rq, _ = h.Repository.GetCalculationByID(id)
	jsonResponse(ctx, rq, 1, gin.H{"id": id})
}

// ApiDeleteRequest удалить заявку
// @Summary Удалить заявку
// @Description Требуется авторизация. Доступно только создателю заявки.
// @Tags requests
// @Security BearerAuth
// @Security CookieAuth
// @Accept json
// @Produce json
// @Param id path int true "ID заявки"
// @Success 200 {object} object{data=object{deleted=uint},total=int,filters=object}
// @Failure 400 {object} object{error=string}
// @Failure 401 {object} object{error=string}
// @Failure 403 {object} object{error=string}
// @Router /api/requests/{id} [delete]
func (h *Handler) ApiDeleteCalculation(ctx *gin.Context) {
	userID, ok := middleware.GetCurrentUserID(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	// только создатель или модератор
	owner, err := h.Repository.IsCalculationOwner(userID, uint(id))
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}
	isModerator := false
	if v, ok := ctx.Get(middleware.IsModeratorKey); ok {
		if b, ok := v.(bool); ok {
			isModerator = b
		}
	}
	if !owner && !isModerator {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}
	if err := h.Repository.DeleteCalculation(uint(id)); err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}
	jsonResponse(ctx, gin.H{"deleted": id}, 1, gin.H{})
}

package handler

import (
	"bytes"
	"io"
	"net/http"

	"partitionlab/internal/app/middleware"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type reqSymptomKey struct {
	// RequestID опционален для добавления: если не указан, используем/создаём черновик
	RequestID   uint `json:"calculation_id"`
	PartitionID uint `json:"partition_id"` // убрал binding:"required" чтобы проверить вручную
}

// ApiAddRequestSymptom добавить симптом в заявку
// @Summary Добавить симптом в заявку
// @Description Требуется авторизация. Добавляет симптом в черновик заявки. Если calculation_id не указан, создается/используется текущий черновик.
// @Tags request-partitions
// @Security BearerAuth
// @Security CookieAuth
// @Accept json
// @Produce json
// @Param request_symptom body object{calculation_id=uint,partition_id=uint} true "Данные для добавления"
// @Success 200 {object} object{data=object{added=uint},total=int,filters=object}
// @Failure 400 {object} object{error=string}
// @Failure 401 {object} object{error=string}
// @Failure 403 {object} object{error=string}
// @Router /api/request-partitions [post]
func (h *Handler) ApiAddCalculationPartition(ctx *gin.Context) {
	userID, ok := middleware.GetCurrentUserID(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var body reqSymptomKey

	// DEBUG: логируем сырое тело запроса
	rawBody, _ := ctx.GetRawData()
	ctx.Request.Body = io.NopCloser(bytes.NewBuffer(rawBody))
	logrus.Infof("DEBUG: Raw body: %s", string(rawBody))

	if err := ctx.ShouldBindJSON(&body); err != nil {
		logrus.Errorf("DEBUG: Bind error: %v, Body: %+v", err, body)
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	logrus.Infof("DEBUG: Parsed body: calculation_id=%d, partition_id=%d", body.RequestID, body.PartitionID)

	// Проверяем partition_id вручную
	if body.PartitionID == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "partition_id is required and must be > 0"})
		return
	}

	var requestID uint
	if body.RequestID == 0 {
		// Автосоздание/получение черновика
		draft, err := h.Repository.GetOrCreateDraftCalculation(userID)
		if err != nil {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
			return
		}
		requestID = draft.ID
	} else {
		requestID = body.RequestID
		// только владелец черновика
		if owner, err := h.Repository.IsCalculationOwner(userID, requestID); err != nil || !owner {
			h.errorHandler(ctx, http.StatusForbidden, err)
			return
		}
		rq, err := h.Repository.GetCalculationByID(requestID)
		if err != nil || rq.Status != "черновик" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "only draft can be modified"})
			return
		}
	}

	if err := h.Repository.AddPartitionToCalculation(requestID, body.PartitionID); err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}
	jsonResponse(ctx, gin.H{"added": body.PartitionID}, 1, gin.H{"calculation_id": requestID})
}

// ApiDeleteRequestSymptom удалить симптом из заявки
// @Summary Удалить симптом из заявки
// @Description Требуется авторизация. Удаляет симптом из черновика заявки. Доступно только владельцу.
// @Tags request-partitions
// @Security BearerAuth
// @Security CookieAuth
// @Accept json
// @Produce json
// @Param request_symptom body object{calculation_id=uint,partition_id=uint} true "Данные для удаления"
// @Success 200 {object} object{data=object{deleted=object},total=int,filters=object}
// @Failure 400 {object} object{error=string}
// @Failure 401 {object} object{error=string}
// @Failure 403 {object} object{error=string}
// @Router /api/request-partitions [delete]
func (h *Handler) ApiDeleteCalculationPartition(ctx *gin.Context) {
	userID, ok := middleware.GetCurrentUserID(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var body reqSymptomKey
	if err := ctx.ShouldBindJSON(&body); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	// только владелец черновика может менять
	if owner, err := h.Repository.IsCalculationOwner(userID, body.RequestID); err != nil || !owner {
		h.errorHandler(ctx, http.StatusForbidden, err)
		return
	}
	rq, err := h.Repository.GetCalculationByID(body.RequestID)
	if err != nil || rq.Status != "черновик" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "only draft can be modified"})
		return
	}
	if err := h.Repository.DeleteCalculationItem(body.RequestID, body.PartitionID); err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}
	jsonResponse(ctx, gin.H{"deleted": body}, 1, gin.H{})
}

// ApiUpdateRequestSymptom обновить данные симптома в заявке
// @Summary Обновить данные симптома в заявке
// @Description Требуется авторизация. Обновляет данные связи заявка-симптом (интенсивность, комментарий, признак основного). Доступно только владельцу черновика.
// @Tags request-partitions
// @Security BearerAuth
// @Security CookieAuth
// @Accept json
// @Produce json
// @Param request_symptom body object{calculation_id=uint,partition_id=uint,quantity=int,comment=string,is_main=bool} true "Данные для обновления"
// @Success 200 {object} object{data=object{updated=uint},total=int,filters=object}
// @Failure 400 {object} object{error=string}
// @Failure 401 {object} object{error=string}
// @Failure 403 {object} object{error=string}
// @Router /api/request-partitions [put]
func (h *Handler) ApiUpdateCalculationPartition(ctx *gin.Context) {
	userID, ok := middleware.GetCurrentUserID(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	type bodyT struct {
		RequestID   uint    `json:"calculation_id" binding:"required"`
		PartitionID uint    `json:"partition_id" binding:"required"`
		Quantity    *int    `json:"quantity"`
		Comment     *string `json:"comment"`
		IsMain      *bool   `json:"is_main"`
	}
	var body bodyT
	if err := ctx.ShouldBindJSON(&body); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	if owner, err := h.Repository.IsCalculationOwner(userID, body.RequestID); err != nil || !owner {
		h.errorHandler(ctx, http.StatusForbidden, err)
		return
	}
	rq, err := h.Repository.GetCalculationByID(body.RequestID)
	if err != nil || rq.Status != "черновик" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "only draft can be modified"})
		return
	}
	if err := h.Repository.UpdateCalculationItem(body.RequestID, body.PartitionID, body.Quantity, body.Comment, body.IsMain); err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}
	jsonResponse(ctx, gin.H{"updated": body.PartitionID}, 1, gin.H{"calculation_id": body.RequestID})
}

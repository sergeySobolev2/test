package handler

import (
	"net/http"
	"strconv"
	"time"

	"partitionlab/internal/app/ds"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetSymptom(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	Partition, err := h.Repository.GetPartition(uint(id))
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"Partition":  Partition,
		"minio_base": h.GetMinIOBaseURL(),
	})
}

func (h *Handler) GetPartitions(ctx *gin.Context) {
	searchQuery := ctx.Query("query")
	var partitions []ds.Partition
	var err error

	if searchQuery == "" {
		partitions, err = h.Repository.GetActivePartitions()
	} else {
		partitions, err = h.Repository.SearchPartitions(searchQuery)
	}

	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	// Count items in user's draft request
	userID := uint(1)
	calculationCount, _ := h.Repository.CountDraftItems(userID)

	ctx.JSON(http.StatusOK, gin.H{
		"time":             time.Now().Format("15:04:05"),
		"partitions":       partitions,
		"query":            searchQuery,
		"calculationCount": calculationCount,
		"minio_base":       h.GetMinIOBaseURL(),
	})
}

func (h *Handler) GetCalculationPage(ctx *gin.Context) {
	userID := uint(1)
	// Возможность открыть конкретную заявку по id: /partition-calc?id=123
	if idStr := ctx.Query("id"); idStr != "" {
		if id, err := strconv.ParseUint(idStr, 10, 64); err == nil {
			// Проверим принадлежность заявки пользователю
			if owner, err := h.Repository.IsCalculationOwner(userID, uint(id)); err == nil && owner {
				req, err := h.Repository.GetCalculationByID(uint(id))
				if err == nil && req != nil {
					syms, _ := h.Repository.GetCalculationPartitions(uint(id))
					ctx.JSON(http.StatusOK, gin.H{
						"partitions":    syms,
						"symptomsCount": len(syms),
						"request":       req,
						"time":          time.Now().Format("15:04:05"),
						"minio_base":    h.GetMinIOBaseURL(),
					})
					return
				}
			}
			// если что-то пошло не так — 404
			ctx.Status(http.StatusNotFound)
			return
		}
	}

	// По умолчанию — работаем с черновиком
	calculationSymptoms, request, err := h.Repository.GetDraftPartitions(userID)
	if err != nil {
		// Если черновика нет - редирект на главную
		ctx.Redirect(http.StatusFound, "/")
		return
	}
	calculationCount := len(calculationSymptoms)
	if calculationCount == 0 {
		// Удаляем пустой черновик чтобы не копился мусор
		_ = h.Repository.DeleteCalculation(request.ID)
		ctx.Redirect(http.StatusFound, "/")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"partitions":    calculationSymptoms,
		"symptomsCount": calculationCount,
		"request":       request,
		"time":          time.Now().Format("15:04:05"),
		"minio_base":    h.GetMinIOBaseURL(),
	})
}

func (h *Handler) AddToRequest(ctx *gin.Context) {
	partitionIDStr := ctx.Param("id")
	partitionID, err := strconv.ParseUint(partitionIDStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Partition ID"})
		return
	}

	userID := uint(1)

	request, err := h.Repository.GetOrCreateDraftCalculation(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}

	err = h.Repository.AddPartitionToCalculation(request.ID, uint(partitionID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add Partition to request"})
		return
	}

	ctx.Redirect(http.StatusFound, "/")
}

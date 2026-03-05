package handler

import (
	"fmt"
	"partitionlab/internal/app/config"
	"partitionlab/internal/app/middleware"
	"partitionlab/internal/app/pkg/auth"
	"partitionlab/internal/app/repository"
	"sync"

	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	Repository     *repository.Repository
	Config         *config.Config
	Storage        interface{}
	JWTService     *auth.JWTService
	SessionService *auth.SessionService
	CurrentUserID  uint
}

func NewHandler(r *repository.Repository, cfg *config.Config, jwtSvc *auth.JWTService, sessSvc *auth.SessionService) *Handler {
	return &Handler{
		Repository:     r,
		Config:         cfg,
		JWTService:     jwtSvc,
		SessionService: sessSvc,
		CurrentUserID:  1,
	}
}

// singleton текущего пользователя (лабораторная: фиксированный пользователь)
// Используем sync.Once, чтобы явно реализовать семантику синглтона функции
// и иметь единое место для возможного расширения (например, вытягивать из контекста).
// По методичке — фиксированный пользователь с id=1.
// В коде ниже обращаемся к CurrentUserID() вместо поля структуры.
var (
	userOnce     sync.Once
	cachedUserID uint
)

// CurrentUserID возвращает идентификатор текущего пользователя как функц-синглтон.
func CurrentUserID() uint {
	userOnce.Do(func() {
		cachedUserID = 1
	})
	return cachedUserID
}

// RegisterHandler Функция, в которой мы отдельно регистрируем маршруты
// func (h *Handler) RegisterHandler(router *gin.Engine) {
// 	router.GET("/", h.GetPartitions)
// 	router.GET("/order/:id", h.GetPartition)
// 	router.GET("/calculation", h.GetCalculationPage)
// 	router.POST("/add-to-request/:id", h.AddToRequest)
// 	router.POST("/delete-request", h.DeleteCalculation)
// 	router.GET("/calculations", h.ListRequests)
// 	router.GET("/request/:id", h.ViewRequest)
// }

// RegisterHandler Функция, в которой мы отдельно регистрируем маршруты
func (h *Handler) RegisterHandler(router *gin.Engine) {
	router.GET("/", h.GetPartitions)
	router.GET("/partition/:id", h.GetSymptom)
	router.GET("/soundproofing-calc", h.GetCalculationPage)
	router.POST("/add-partition/:id", h.AddToRequest)
	router.POST("/clear-calc", h.DeleteRequest)
	router.GET("/calculation-history", h.ListRequests)
	router.GET("/calculation-case/:id", h.ViewRequest)
}

// GetMinIOBaseURL возвращает базовый URL для MinIO
func (h *Handler) GetMinIOBaseURL() string {
	return fmt.Sprintf("http://%s:%s", h.Config.MinIOHost, h.Config.MinIOPort)
}

// BuildPublicImageURL собирает публичный URL для изображения из ключа (image_url)
func (h *Handler) BuildPublicImageURL(key string) string {
	if key == "" {
		return ""
	}
	bucket := url.PathEscape(h.Config.MinIOBucket)
	escapedKey := url.PathEscape(key)
	return fmt.Sprintf("%s/%s/%s", h.GetMinIOBaseURL(), bucket, escapedKey)
}

// errorHandler для более удобного вывода ошибок
func (h *Handler) errorHandler(ctx *gin.Context, errorStatusCode int, err error) {
	msg := "unknown error"
	if err != nil {
		msg = err.Error()
	}
	logrus.Error(msg)
	ctx.JSON(errorStatusCode, gin.H{
		"status":      "error",
		"description": msg,
	})
}

// RegisterAPI регистрирует REST API маршруты
func (h *Handler) RegisterAPI(router *gin.Engine) {
	api := router.Group("/api")

	// Public auth endpoints (без авторизации)
	api.POST("/users/register", h.ApiRegisterUser)
	api.POST("/users/login", h.ApiLogin)

	// Импортируем middleware для удобства
	authSvc := &middleware.AuthService{
		JWT:     h.JWTService,
		Session: h.SessionService,
	}

	// Public read endpoints (без авторизации, опциональная аутентификация для фильтрации)
	publicSymptoms := api.Group("/partitions")
	{
		publicSymptoms.GET("", h.ApiListPartitions)
		publicSymptoms.GET("/:id", h.ApiGetPartition)
	}

	// Protected endpoints (требуют аутентификации)
	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware(authSvc))
	{
		// User profile
		protected.POST("/users/logout", h.ApiLogout)
		protected.GET("/users/profile", h.ApiGetProfile)
		protected.PUT("/users/profile", h.ApiUpdateProfile)

		// Requests (требуют авторизации)
		protected.GET("/requests/cart", h.ApiGetCart)
		protected.GET("/calculations", h.ApiListCalculations)
		protected.GET("/requests/:id", h.ApiGetCalculation)
		protected.PUT("/requests/:id", h.ApiUpdateCalculation)
		protected.PUT("/requests/:id/form", h.ApiFormCalculation)
		protected.DELETE("/requests/:id", h.ApiDeleteCalculation)

		// Request-partitions (требуют авторизации)
		protected.POST("/request-partitions", h.ApiAddCalculationPartition)
		protected.DELETE("/request-partitions", h.ApiDeleteCalculationPartition)
		protected.PUT("/request-partitions", h.ApiUpdateCalculationPartition)
	}

	// Moderator endpoints (требуют роль модератора)
	moderator := api.Group("")
	moderator.Use(middleware.AuthMiddleware(authSvc))
	moderator.Use(middleware.RequireModeratorMiddleware())
	{
		// partitions (CRUD для модератора)
		moderator.POST("/partitions", h.ApiCreatePartition)
		moderator.PUT("/partitions/:id", h.ApiUpdatePartition)
		moderator.DELETE("/partitions/:id", h.ApiDeletePartition)
		moderator.POST("/partitions/:id/image", h.ApiUploadPartitionImage)

		// Complete/reject requests (только модератор)
		moderator.PUT("/requests/:id/complete", h.ApiCompleteCalculation)

		// Trigger async calculation (только модератор)
		moderator.POST("/requests/:id/calculate", h.ApiTriggerAsyncCalculation)
	}

	// Async service endpoints (без middleware - только проверка токена)
	async := api.Group("/async")
	{
		// Receive result from async service
		async.POST("/result", h.ApiReceiveAsyncResult)
	}
}

// jsonResponse — единый формат ответа
func jsonResponse(ctx *gin.Context, data interface{}, total int64, filters gin.H) {
	ctx.JSON(200, gin.H{
		"data":    data,
		"total":   total,
		"filters": filters,
	})
}

// func calculatePartitions(weight float64, partitions []Partition) float64 {
//     // Расчет процента обезвоживания на основе симптомов
//     totalSeverity := 0.0
//     symptomCount := 0

//     // Суммируем тяжесть симптомов
//     for _, Partition := range partitions {
//         switch Partition.Severity {
//         case "Легкая (1-2%)":
//             totalSeverity += 1.5
//         case "Средняя (3-6%)":
//             totalSeverity += 4.5
//         case "Тяжелая (7-9%)":
//             totalSeverity += 8.0
//         }
//         symptomCount++
//     }

//     // Рассчитываем средний процент обезвоживания
//     if symptomCount > 0 {
//         NoiseReductionDB := totalSeverity / float64(symptomCount)
//         return NoiseReductionDB
//     }

//     // Если симптомов нет, возвращаем 0
//     return 0.0
// }

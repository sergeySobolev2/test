package main

import (
	"fmt"

	_ "partitionlab/docs" // Swagger docs
	"partitionlab/internal/app/config"
	"partitionlab/internal/app/dsn"
	"partitionlab/internal/app/handler"
	"partitionlab/internal/app/middleware"
	"partitionlab/internal/app/pkg/auth"
	"partitionlab/internal/app/pkg/storage"
	"partitionlab/internal/app/repository"
	"partitionlab/internal/pkg"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

// @title Partition Lab API
// @version 4.0
// @description API для управления заявками на расчет звукоизоляции перегородок
// @description
// @description Без авторизации доступны:
// @description - Регистрация и аутентификация
// @description - Просмотр симптомов (GET)
// @description
// @description С авторизацией (пользователь):
// @description - Управление своими заявками
// @description - Добавление симптомов в заявки
// @description
// @description С ролью модератора:
// @description - Все методы пользователя
// @description - CRUD операции с симптомами
// @description - Завершение/отклонение заявок

// @contact.name API Support
// @contact.email support@partitionlab.com

// @host localhost:8082
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Введите JWT токен в формате: Bearer {token}

// @securityDefinitions.apikey CookieAuth
// @in header
// @name Cookie
// @description Сессия через cookie (session_id)

func main() {
	_ = godotenv.Load("../../.env")

	router := gin.Default()
	router.Use(middleware.CORSMiddleware())
	conf, err := config.NewConfig()
	if err != nil {
		logrus.Fatalf("error loading config: %v", err)
	}

	postgresString := dsn.FromEnv()
	fmt.Println(postgresString)

	rep, errRep := repository.New(postgresString)
	if errRep != nil {
		logrus.Fatalf("error initializing repository: %v", errRep)
	}

	// Init JWT Service
	jwtService := auth.NewJWTService(conf.JWTSecret)
	logrus.Info("JWT service initialized")

	// Init Session Service (Redis)
	sessionService, err := auth.NewSessionService(conf.RedisHost, conf.RedisPort, conf.RedisPassword, conf.RedisDB)
	if err != nil {
		logrus.Fatalf("error initializing session service: %v", err)
	}
	logrus.Info("Session service (Redis) initialized")

	hand := handler.NewHandler(rep, conf, jwtService, sessionService)

	// Init MinIO storage
	publicBase := fmt.Sprintf("http://%s:%s", conf.MinIOHost, conf.MinIOPort)
	minioClient, err := storage.NewMinIO(conf.MinIOHost+":"+conf.MinIOPort, conf.MinIOAccessKey, conf.MinIOSecretKey, conf.MinIOBucket, conf.MinIOUseSSL, publicBase)
	if err != nil {
		logrus.Warnf("minio init failed: %v", err)
	} else {
		hand.Storage = minioClient
	}

	application := pkg.NewApp(conf, router, hand)
	application.RunApp()
}

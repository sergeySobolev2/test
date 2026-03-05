package pkg

import (
	"fmt"
	"os"
	"strconv"

	"partitionlab/internal/app/config"
	"partitionlab/internal/app/handler"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Application struct {
	Config  *config.Config
	Router  *gin.Engine
	Handler *handler.Handler
}

func NewApp(c *config.Config, r *gin.Engine, h *handler.Handler) *Application {
	return &Application{
		Config:  c,
		Router:  r,
		Handler: h,
	}
}

func (a *Application) RunApp() {
	logrus.Info("Server start up")

	// Swagger UI
	a.Router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	a.Handler.RegisterHandler(a.Router)
	a.Handler.RegisterAPI(a.Router)

	port := a.Config.ServicePort
	if envPort := os.Getenv("PORT"); envPort != "" {
		if parsedPort, err := strconv.Atoi(envPort); err == nil {
			port = parsedPort
		}
	}

	serverAddress := fmt.Sprintf("%s:%d", a.Config.ServiceHost, port)

	useTLS := os.Getenv("ENABLE_TLS") == "true"
	certFile := os.Getenv("TLS_CERT_FILE")
	keyFile := os.Getenv("TLS_KEY_FILE")
	if useTLS && certFile != "" && keyFile != "" {
		if _, errCert := os.Stat(certFile); errCert == nil {
			if _, errKey := os.Stat(keyFile); errKey == nil {
				logrus.Infof("HTTPS enabled: %s", serverAddress)
				if err := a.Router.RunTLS(serverAddress, certFile, keyFile); err != nil {
					logrus.Fatal(err)
				}
				return
			}
		}
		logrus.Warn("TLS enabled but cert/key not found. Falling back to HTTP.")
	}

	if err := a.Router.Run(serverAddress); err != nil {
		logrus.Fatal(err)
	}
	logrus.Info("Server down")
}

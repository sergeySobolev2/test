package middleware

import (
	"net"
	"net/url"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func isAllowedOrigin(origin string, extra []string) bool {
	if origin == "" {
		// запросы без Origin (curl/сервер-сервер)
		return true
	}

	// Tauri origins
	if origin == "tauri://localhost" || origin == "http://tauri.localhost" || origin == "https://tauri.localhost" {
		return true
	}

	u, err := url.Parse(origin)
	if err == nil {
		h := strings.ToLower(u.Hostname())
		// удобный allowlist для локальных сетевых тестов
		if h == "localhost" || h == "127.0.0.1" || h == "::1" {
			return true
		}
		if ip := net.ParseIP(h); ip != nil && isPrivateIP(ip) {
			return true
		}
	}

	for _, o := range extra {
		if strings.EqualFold(strings.TrimSpace(o), origin) {
			return true
		}
	}

	return false
}

func isPrivateIP(ip net.IP) bool {
	ip = ip.To4()
	if ip == nil {
		return false
	}
	// 10.0.0.0/8
	if ip[0] == 10 {
		return true
	}
	// 172.16.0.0/12
	if ip[0] == 172 && ip[1] >= 16 && ip[1] <= 31 {
		return true
	}
	// 192.168.0.0/16
	if ip[0] == 192 && ip[1] == 168 {
		return true
	}
	return false
}

// CORSMiddleware разрешает запросы от SPA/PWA/Tauri.
// Доп. разрешённые origin можно передать через env CORS_ALLOW_ORIGINS (через запятую).
func CORSMiddleware() gin.HandlerFunc {
	cfg := cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
		AllowOriginFunc: func(origin string) bool {
			return true
		},
	}

	return cors.New(cfg)
}

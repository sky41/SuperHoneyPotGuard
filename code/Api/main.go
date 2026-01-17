package main

import (
	"log"
	"superhoneypotguard/config"
	"superhoneypotguard/database"
	"superhoneypotguard/middleware"
	"superhoneypotguard/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadConfig()

	cfg := config.AppConfig
	if cfg.GinMode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	database.InitDB()

	middleware.InitRateLimiter()

	r := gin.Default()

	r.Use(middleware.RateLimitMiddleware())
	r.Use(middleware.LogMiddleware())

	routes.SetupRoutes(r)

	addr := ":" + cfg.Port
	log.Printf("服务器运行在 http://localhost%s", addr)
	log.Printf("环境: %s", cfg.GinMode)

	if err := r.Run(addr); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}

package main

import (
	"context"
	"log"

	"go_template/internal/app"
	"go_template/internal/config"
	"go_template/internal/health"

	"github.com/gin-gonic/gin"
)

func main() {
	ctx := context.Background()
	application, err := app.Bootstrap(ctx, config.Options{Runtime: config.RuntimeAPI})
	if err != nil {
		log.Fatalf("bootstrap api: %v", err)
	}
	defer application.Shutdown(ctx)

	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())
	health.RegisterRoutes(router, nil)

	if err := router.Run(":" + application.Config.Server.Port); err != nil {
		log.Fatalf("run api server: %v", err)
	}
}

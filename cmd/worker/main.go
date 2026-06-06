package main

import (
	"context"
	"log"

	"go_template/internal/app"
	"go_template/internal/config"
)

func main() {
	ctx := context.Background()
	application, err := app.Bootstrap(ctx, config.Options{Runtime: config.RuntimeWorker})
	if err != nil {
		log.Fatalf("bootstrap worker: %v", err)
	}
	defer application.Shutdown(ctx)

	application.Logger.Info("worker bootstrap complete; queue consumption remains on cmd/queue during template transition")
}

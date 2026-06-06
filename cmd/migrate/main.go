package main

import (
	"context"
	"log"

	"go_template/internal/app"
	"go_template/internal/config"
)

func main() {
	ctx := context.Background()
	application, err := app.Bootstrap(ctx, config.Options{Runtime: config.RuntimeMigrate})
	if err != nil {
		log.Fatalf("bootstrap migrate: %v", err)
	}
	defer application.Shutdown(ctx)

	application.Logger.Info("migration bootstrap complete; schema execution remains on existing migration workflow during template transition")
}

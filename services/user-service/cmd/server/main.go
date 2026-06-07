package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"user_service/internal/httpapi"
	"user_service/internal/user"
)

func main() {
	serviceToken := os.Getenv("SERVICE_TOKEN")
	if serviceToken == "" {
		log.Fatal("SERVICE_TOKEN is required")
	}

	router := httpapi.NewRouter(user.NewMemoryRepository(time.Now), serviceToken)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	server := &http.Server{Addr: ":" + port, Handler: router, ReadHeaderTimeout: 5 * time.Second}
	log.Printf("user-service listening on :%s", port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

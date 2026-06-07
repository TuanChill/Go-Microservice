package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"notification_otp_service/internal/httpapi"
	"notification_otp_service/internal/notification"
	"notification_otp_service/internal/otp"
)

func main() {
	repo := otp.NewMemoryRepository()
	otpService := otp.NewService(repo, otp.GenerateSixDigitCode, time.Now, 5*time.Minute)
	serviceToken := os.Getenv("SERVICE_TOKEN")
	if serviceToken == "" {
		log.Fatal("SERVICE_TOKEN is required")
	}
	router := httpapi.NewRouter(otpService, notification.NewService(), serviceToken)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	server := &http.Server{Addr: ":" + port, Handler: router, ReadHeaderTimeout: 5 * time.Second}
	log.Printf("notification-otp-service listening on :%s", port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

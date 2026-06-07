package httpapi

import (
	"errors"
	"net/http"
	"strings"

	"notification_otp_service/internal/notification"
	"notification_otp_service/internal/otp"

	"github.com/gin-gonic/gin"
)

type Server struct {
	otp          *otp.Service
	notification *notification.Service
	serviceToken string
}

func NewRouter(otpService *otp.Service, notificationService *notification.Service, serviceToken string) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	server := &Server{otp: otpService, notification: notificationService, serviceToken: serviceToken}

	internal := router.Group("/internal", server.requireServiceHeaders)
	internal.POST("/otp/request", server.requestOtp)
	internal.POST("/otp/verify", server.verifyOtp)
	internal.POST("/notifications/email-verification", server.sendEmailVerification)
	internal.POST("/notifications/password-reset", server.sendPasswordReset)
	return router
}

func (s *Server) requireServiceHeaders(c *gin.Context) {
	authorization := c.GetHeader("Authorization")
	if !strings.HasPrefix(authorization, "Bearer ") || strings.TrimSpace(strings.TrimPrefix(authorization, "Bearer ")) != s.serviceToken {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	if len(c.GetHeader("X-Correlation-ID")) < 8 || len(c.GetHeader("Idempotency-Key")) < 16 {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	c.Next()
}

func (s *Server) requestOtp(c *gin.Context) {
	var req otp.Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	res, err := s.otp.Request(req)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusCreated, res)
}

func (s *Server) verifyOtp(c *gin.Context) {
	var req otp.VerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	res, err := s.otp.Verify(req)
	if errors.Is(err, otp.ErrNotFound) {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, res)
}

func (s *Server) sendEmailVerification(c *gin.Context) {
	var req notification.EmailVerificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	c.JSON(http.StatusAccepted, s.notification.AcceptEmailVerification(c.GetHeader("Idempotency-Key"), req))
}

func (s *Server) sendPasswordReset(c *gin.Context) {
	var req notification.PasswordResetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	c.JSON(http.StatusAccepted, s.notification.AcceptPasswordReset(c.GetHeader("Idempotency-Key"), req))
}

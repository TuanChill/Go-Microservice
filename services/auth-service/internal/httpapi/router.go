package httpapi

import (
	"errors"
	"net/http"
	"strings"

	"auth_service/internal/auth"

	"github.com/gin-gonic/gin"
)

type Server struct {
	auth         *auth.Service
	serviceToken string
}

func NewRouter(authService *auth.Service, serviceToken string) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	server := &Server{auth: authService, serviceToken: serviceToken}

	internal := router.Group("/internal", server.requireServiceHeaders)
	internal.POST("/auth/register", server.register)
	internal.POST("/auth/login", server.login)
	internal.POST("/auth/refresh", server.refresh)
	internal.POST("/auth/verify-account", server.verifyAccount)
	internal.POST("/auth/password-reset", server.resetPassword)
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

func (s *Server) register(c *gin.Context) {
	var req auth.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	res, err := s.auth.Register(req)
	if errors.Is(err, auth.ErrConflict) {
		c.AbortWithStatus(http.StatusConflict)
		return
	}
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusCreated, res)
}

func (s *Server) login(c *gin.Context) {
	var req auth.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	res, err := s.auth.Login(req)
	if errors.Is(err, auth.ErrInvalidCredentials) {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, res)
}

func (s *Server) refresh(c *gin.Context) {
	var req auth.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	res, err := s.auth.Refresh(req)
	if errors.Is(err, auth.ErrInvalidToken) {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, res)
}

func (s *Server) verifyAccount(c *gin.Context) {
	var req auth.VerifyAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	res, err := s.auth.VerifyAccount(req)
	if errors.Is(err, auth.ErrInvalidToken) {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, res)
}

func (s *Server) resetPassword(c *gin.Context) {
	var req auth.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	res, err := s.auth.ResetPassword(req)
	if errors.Is(err, auth.ErrInvalidToken) {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, res)
}

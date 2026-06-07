package httpapi

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"user_service/internal/user"

	"github.com/gin-gonic/gin"
)

type Server struct {
	repo         user.Repository
	serviceToken string
}

func NewRouter(repo user.Repository, serviceToken string) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	server := &Server{repo: repo, serviceToken: serviceToken}

	internal := router.Group("/internal", server.requireCorrelationAndAuth)
	internal.POST("/users", server.requireIdempotency, server.createUser)
	internal.GET("/users/:id", server.getUserProfile)
	internal.PATCH("/users/:id", server.requireIdempotency, server.updateUserProfile)
	internal.DELETE("/users/:id", server.requireIdempotency, server.deactivateUser)
	return router
}

func (s *Server) requireCorrelationAndAuth(c *gin.Context) {
	authorization := c.GetHeader("Authorization")
	if !strings.HasPrefix(authorization, "Bearer ") || strings.TrimSpace(strings.TrimPrefix(authorization, "Bearer ")) != s.serviceToken {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	if len(c.GetHeader("X-Correlation-ID")) < 8 {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	c.Next()
}

func (s *Server) requireIdempotency(c *gin.Context) {
	if len(c.GetHeader("Idempotency-Key")) < 16 {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	c.Next()
}

func (s *Server) createUser(c *gin.Context) {
	var req user.CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	res, err := s.repo.Create(req.Email)
	if errors.Is(err, user.ErrConflict) {
		c.AbortWithStatus(http.StatusConflict)
		return
	}
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusCreated, res)
}

func (s *Server) getUserProfile(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	profile, err := s.repo.GetProfile(id)
	if errors.Is(err, user.ErrNotFound) {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, profile)
}

func (s *Server) updateUserProfile(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var req user.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	profile, err := s.repo.UpdateProfile(id, req)
	if errors.Is(err, user.ErrNotFound) {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, profile)
}

func (s *Server) deactivateUser(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	res, err := s.repo.Deactivate(id)
	if errors.Is(err, user.ErrNotFound) {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, res)
}

func parseID(c *gin.Context) (int, bool) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id < 1 {
		c.AbortWithStatus(http.StatusBadRequest)
		return 0, false
	}
	return id, true
}

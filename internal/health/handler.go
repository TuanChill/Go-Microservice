package health

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, readiness func() error) {
	router.GET("/health/live", live)
	router.GET("/health/ready", ready(readiness))
}

func live(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func ready(readiness func() error) gin.HandlerFunc {
	return func(c *gin.Context) {
		if readiness != nil {
			if err := readiness(); err != nil {
				c.JSON(http.StatusServiceUnavailable, gin.H{"status": "unavailable"})
				return
			}
		}
		c.JSON(http.StatusOK, gin.H{"status": "ready"})
	}
}

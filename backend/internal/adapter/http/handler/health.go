package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type HealthHandler struct {
	systemDB *gorm.DB
}

func NewHealthHandler(systemDB *gorm.DB) *HealthHandler {
	return &HealthHandler{systemDB: systemDB}
}

func (h *HealthHandler) Healthz(c *gin.Context) {
	sqlDB, err := h.systemDB.DB()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"status": "unhealthy", "error": "db connection error"})
		return
	}
	if err := sqlDB.Ping(); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"status": "unhealthy", "error": "db ping failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "healthy"})
}

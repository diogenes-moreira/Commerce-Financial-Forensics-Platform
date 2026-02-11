package handler

import (
	"net/http"

	"github.com/diogenes/costforensics/backend/internal/adapter/http/dto/response"
	tenantrepo "github.com/diogenes/costforensics/backend/internal/adapter/postgres/tenant"
	anomalyapp "github.com/diogenes/costforensics/backend/internal/application/anomaly"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AnomalyHandler struct{}

func NewAnomalyHandler() *AnomalyHandler {
	return &AnomalyHandler{}
}

func (h *AnomalyHandler) getService(c *gin.Context) *anomalyapp.Service {
	db := c.MustGet("tenant_db").(*gorm.DB)
	repo := tenantrepo.NewAnomalyRepo(db)
	return anomalyapp.NewService(repo)
}

func (h *AnomalyHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	svc := h.getService(c)
	a, err := svc.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, response.AnomalyFromDomain(a))
}

func (h *AnomalyHandler) List(c *gin.Context) {
	svc := h.getService(c)
	anomalies, err := svc.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	items := make([]response.Anomaly, len(anomalies))
	for i, a := range anomalies {
		items[i] = response.AnomalyFromDomain(a)
	}
	c.JSON(http.StatusOK, items)
}

func (h *AnomalyHandler) Resolve(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	svc := h.getService(c)
	a, err := svc.Resolve(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, response.AnomalyFromDomain(a))
}

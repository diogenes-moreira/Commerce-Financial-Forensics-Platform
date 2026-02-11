package handler

import (
	"io"
	"net/http"

	"github.com/diogenes/costforensics/backend/internal/adapter/http/dto/request"
	"github.com/diogenes/costforensics/backend/internal/adapter/http/dto/response"
	tenantrepo "github.com/diogenes/costforensics/backend/internal/adapter/postgres/tenant"
	integrationapp "github.com/diogenes/costforensics/backend/internal/application/integration"
	"github.com/diogenes/costforensics/backend/internal/domain/integration"
	"github.com/diogenes/costforensics/backend/internal/domain/shared"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IntegrationHandler struct{}

func NewIntegrationHandler() *IntegrationHandler {
	return &IntegrationHandler{}
}

func (h *IntegrationHandler) getService(c *gin.Context) *integrationapp.Service {
	db := c.MustGet("tenant_db").(*gorm.DB)
	repo := tenantrepo.NewIntegrationRepo(db)
	return integrationapp.NewService(repo)
}

func (h *IntegrationHandler) Create(c *gin.Context) {
	var req request.CreateIntegration
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	svc := h.getService(c)
	i, err := svc.Create(c.Request.Context(), req.Platform, req.Name, req.Config)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, response.IntegrationFromDomain(i))
}

func (h *IntegrationHandler) List(c *gin.Context) {
	var req request.ListIntegrations
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	filter := integration.ListFilter{
		Platform:   req.Platform,
		Status:     integration.IntegrationStatus(req.Status),
		Pagination: shared.NewPagination(req.Page, req.PageSize),
	}

	svc := h.getService(c)
	result, err := svc.List(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	items := make([]response.Integration, len(result.Items))
	for idx, i := range result.Items {
		items[idx] = response.IntegrationFromDomain(i)
	}
	c.JSON(http.StatusOK, response.Paginated[response.Integration]{
		Items: items, TotalCount: result.TotalCount,
		Page: result.Page, PageSize: result.PageSize,
	})
}

func (h *IntegrationHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	svc := h.getService(c)
	i, err := svc.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, response.IntegrationFromDomain(i))
}

func (h *IntegrationHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req request.UpdateIntegration
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	svc := h.getService(c)
	i, err := svc.Update(c.Request.Context(), id, req.Name, req.Config)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, response.IntegrationFromDomain(i))
}

func (h *IntegrationHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	svc := h.getService(c)
	if err := svc.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

func (h *IntegrationHandler) TriggerSync(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	svc := h.getService(c)
	i, err := svc.TriggerSync(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, response.IntegrationFromDomain(i))
}

func (h *IntegrationHandler) HandleWebhook(c *gin.Context) {
	integrationID, err := uuid.Parse(c.Param("integration_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid integration_id"})
		return
	}

	svc := h.getService(c)
	i, err := svc.GetByID(c.Request.Context(), integrationID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	conn, err := svc.GetConnector(i)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	eventType := c.GetHeader("X-Webhook-Event")
	payload, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read request body"})
		return
	}

	if err := conn.HandleWebhook(c.Request.Context(), eventType, payload); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

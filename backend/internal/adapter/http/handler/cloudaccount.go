package handler

import (
	"net/http"

	"github.com/diogenes/costforensics/backend/internal/adapter/http/dto/request"
	"github.com/diogenes/costforensics/backend/internal/adapter/http/dto/response"
	tenantrepo "github.com/diogenes/costforensics/backend/internal/adapter/postgres/tenant"
	cloudaccountapp "github.com/diogenes/costforensics/backend/internal/application/cloudaccount"
	"github.com/diogenes/costforensics/backend/internal/domain/cloudaccount"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CloudAccountHandler struct{}

func NewCloudAccountHandler() *CloudAccountHandler {
	return &CloudAccountHandler{}
}

func (h *CloudAccountHandler) getService(c *gin.Context) *cloudaccountapp.Service {
	db := c.MustGet("tenant_db").(*gorm.DB)
	repo := tenantrepo.NewCloudAccountRepo(db)
	return cloudaccountapp.NewService(repo)
}

func (h *CloudAccountHandler) Create(c *gin.Context) {
	var req request.CreateCloudAccount
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	svc := h.getService(c)
	a, err := svc.Create(c.Request.Context(), cloudaccount.Provider(req.Provider), req.Name, req.ExternalID)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, response.CloudAccountFromDomain(a))
}

func (h *CloudAccountHandler) GetByID(c *gin.Context) {
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
	c.JSON(http.StatusOK, response.CloudAccountFromDomain(a))
}

func (h *CloudAccountHandler) List(c *gin.Context) {
	svc := h.getService(c)
	accounts, err := svc.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	items := make([]response.CloudAccount, len(accounts))
	for i, a := range accounts {
		items[i] = response.CloudAccountFromDomain(a)
	}
	c.JSON(http.StatusOK, items)
}

func (h *CloudAccountHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req request.UpdateCloudAccount
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	svc := h.getService(c)
	a, err := svc.Update(c.Request.Context(), id, req.Name)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, response.CloudAccountFromDomain(a))
}

func (h *CloudAccountHandler) Delete(c *gin.Context) {
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

func (h *CloudAccountHandler) BeginSync(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	svc := h.getService(c)
	a, err := svc.BeginSync(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, response.CloudAccountFromDomain(a))
}

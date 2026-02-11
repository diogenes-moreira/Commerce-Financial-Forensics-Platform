package handler

import (
	"net/http"

	"github.com/diogenes/costforensics/backend/internal/adapter/http/dto/request"
	"github.com/diogenes/costforensics/backend/internal/adapter/http/dto/response"
	tenantrepo "github.com/diogenes/costforensics/backend/internal/adapter/postgres/tenant"
	sellerapp "github.com/diogenes/costforensics/backend/internal/application/seller"
	"github.com/diogenes/costforensics/backend/internal/domain/seller"
	"github.com/diogenes/costforensics/backend/internal/domain/shared"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SellerHandler struct{}

func NewSellerHandler() *SellerHandler {
	return &SellerHandler{}
}

func (h *SellerHandler) getService(c *gin.Context) *sellerapp.Service {
	db := c.MustGet("tenant_db").(*gorm.DB)
	repo := tenantrepo.NewSellerRepo(db)
	return sellerapp.NewService(repo)
}

func (h *SellerHandler) Create(c *gin.Context) {
	var req request.CreateSeller
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	svc := h.getService(c)
	s, err := svc.Create(c.Request.Context(), req.ExternalID, req.Code, req.Name, req.Email, req.CommissionPct)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, response.SellerFromDomain(s))
}

func (h *SellerHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	svc := h.getService(c)
	s, err := svc.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, response.SellerFromDomain(s))
}

func (h *SellerHandler) List(c *gin.Context) {
	var req request.ListSellers
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	filter := seller.ListFilter{
		Status:     seller.Status(req.Status),
		Pagination: shared.NewPagination(req.Page, req.PageSize),
	}

	svc := h.getService(c)
	result, err := svc.List(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	items := make([]response.Seller, len(result.Items))
	for i, s := range result.Items {
		items[i] = response.SellerFromDomain(s)
	}
	c.JSON(http.StatusOK, response.Paginated[response.Seller]{
		Items: items, TotalCount: result.TotalCount,
		Page: result.Page, PageSize: result.PageSize,
	})
}

func (h *SellerHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req request.UpdateSeller
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	svc := h.getService(c)
	s, err := svc.Update(c.Request.Context(), id, req.Name, req.Email, req.CommissionPct)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, response.SellerFromDomain(s))
}

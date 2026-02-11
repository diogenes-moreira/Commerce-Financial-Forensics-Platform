package handler

import (
	"net/http"

	"github.com/diogenes/costforensics/backend/internal/adapter/http/dto/request"
	"github.com/diogenes/costforensics/backend/internal/adapter/http/dto/response"
	tenantrepo "github.com/diogenes/costforensics/backend/internal/adapter/postgres/tenant"
	productapp "github.com/diogenes/costforensics/backend/internal/application/product"
	"github.com/diogenes/costforensics/backend/internal/domain/product"
	"github.com/diogenes/costforensics/backend/internal/domain/shared"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductHandler struct{}

func NewProductHandler() *ProductHandler {
	return &ProductHandler{}
}

func (h *ProductHandler) getService(c *gin.Context) *productapp.Service {
	db := c.MustGet("tenant_db").(*gorm.DB)
	repo := tenantrepo.NewProductRepo(db)
	return productapp.NewService(repo)
}

func (h *ProductHandler) Create(c *gin.Context) {
	var req request.CreateProduct
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	currency := req.Currency
	if currency == "" {
		currency = "USD"
	}
	unitCost, err := shared.NewMoney(req.UnitCostCents, currency)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	unitPrice, err := shared.NewMoney(req.UnitPriceCents, currency)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	svc := h.getService(c)
	p, err := svc.Create(c.Request.Context(), req.SKU, req.Name, req.Description, req.Category, unitCost, unitPrice)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, response.ProductFromDomain(p))
}

func (h *ProductHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	svc := h.getService(c)
	p, err := svc.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, response.ProductFromDomain(p))
}

func (h *ProductHandler) List(c *gin.Context) {
	var req request.ListProducts
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	filter := product.ListFilter{
		Category:   req.Category,
		Status:     product.Status(req.Status),
		Pagination: shared.NewPagination(req.Page, req.PageSize),
	}

	svc := h.getService(c)
	result, err := svc.List(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	items := make([]response.Product, len(result.Items))
	for i, p := range result.Items {
		items[i] = response.ProductFromDomain(p)
	}
	c.JSON(http.StatusOK, response.Paginated[response.Product]{
		Items: items, TotalCount: result.TotalCount,
		Page: result.Page, PageSize: result.PageSize,
	})
}

func (h *ProductHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req request.UpdateProduct
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var unitPrice *shared.Money
	if req.UnitPriceCents != nil {
		currency := req.Currency
		if currency == "" {
			currency = "USD"
		}
		m, err := shared.NewMoney(*req.UnitPriceCents, currency)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		unitPrice = &m
	}

	svc := h.getService(c)
	p, err := svc.Update(c.Request.Context(), id, req.Name, req.Description, req.Category, unitPrice)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, response.ProductFromDomain(p))
}

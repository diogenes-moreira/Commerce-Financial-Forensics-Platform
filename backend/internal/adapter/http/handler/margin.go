package handler

import (
	"net/http"

	"github.com/diogenes/costforensics/backend/internal/adapter/http/dto/request"
	"github.com/diogenes/costforensics/backend/internal/adapter/http/dto/response"
	tenantrepo "github.com/diogenes/costforensics/backend/internal/adapter/postgres/tenant"
	marginapp "github.com/diogenes/costforensics/backend/internal/application/margin"
	"github.com/diogenes/costforensics/backend/internal/domain/shared"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// MarginHandler handles margin analysis endpoints.
type MarginHandler struct{}

// NewMarginHandler creates a new MarginHandler.
func NewMarginHandler() *MarginHandler {
	return &MarginHandler{}
}

func (h *MarginHandler) getService(c *gin.Context) *marginapp.Service {
	db := c.MustGet("tenant_db").(*gorm.DB)
	orderRepo := tenantrepo.NewOrderRepo(db)
	sellerRepo := tenantrepo.NewSellerRepo(db)
	discountRepo := tenantrepo.NewDiscountApplicationRepo(db)
	return marginapp.NewService(orderRepo, sellerRepo, discountRepo)
}

// GetOrderMargin returns the full margin breakdown for a single order.
func (h *MarginHandler) GetOrderMargin(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	svc := h.getService(c)
	breakdown, err := svc.CalculateForOrder(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response.MarginBreakdownFromDomain(breakdown))
}

// ListMargins returns margin breakdowns for orders within a date range.
func (h *MarginHandler) ListMargins(c *gin.Context) {
	var req request.ListMargins
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	startDate, err := request.ParseDate(req.StartDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start_date format, use YYYY-MM-DD"})
		return
	}
	endDate, err := request.ParseDate(req.EndDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end_date format, use YYYY-MM-DD"})
		return
	}

	var sellerID *uuid.UUID
	if req.SellerID != "" {
		id, err := uuid.Parse(req.SellerID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid seller_id"})
			return
		}
		sellerID = &id
	}

	page := req.Page
	pageSize := req.PageSize
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = shared.DefaultPageSize
	}

	svc := h.getService(c)
	breakdowns, totalCount, err := svc.CalculateForPeriod(c.Request.Context(), startDate, endDate, sellerID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	items := make([]response.MarginBreakdown, len(breakdowns))
	for i, b := range breakdowns {
		items[i] = response.MarginBreakdownFromDomain(&b)
	}

	c.JSON(http.StatusOK, response.Paginated[response.MarginBreakdown]{
		Items:      items,
		TotalCount: totalCount,
		Page:       page,
		PageSize:   pageSize,
	})
}

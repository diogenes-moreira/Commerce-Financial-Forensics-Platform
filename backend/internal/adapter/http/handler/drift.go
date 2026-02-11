package handler

import (
	"net/http"

	"github.com/diogenes/costforensics/backend/internal/adapter/http/dto/request"
	"github.com/diogenes/costforensics/backend/internal/adapter/http/dto/response"
	tenantrepo "github.com/diogenes/costforensics/backend/internal/adapter/postgres/tenant"
	driftapp "github.com/diogenes/costforensics/backend/internal/application/drift"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// DriftHandler handles price/cost drift detection endpoints.
type DriftHandler struct{}

// NewDriftHandler creates a new DriftHandler.
func NewDriftHandler() *DriftHandler {
	return &DriftHandler{}
}

func (h *DriftHandler) getService(c *gin.Context) *driftapp.Service {
	db := c.MustGet("tenant_db").(*gorm.DB)
	orderRepo := tenantrepo.NewOrderRepo(db)
	productRepo := tenantrepo.NewProductRepo(db)
	sellerRepo := tenantrepo.NewSellerRepo(db)
	return driftapp.NewService(orderRepo, productRepo, sellerRepo)
}

// GetOrderDrift returns drift analysis for a single order.
func (h *DriftHandler) GetOrderDrift(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	svc := h.getService(c)
	results, err := svc.DetectForOrder(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	items := make([]response.DriftResult, len(results))
	for i, r := range results {
		items[i] = response.DriftResultFromDomain(r)
	}

	c.JSON(http.StatusOK, gin.H{"items": items})
}

// GetDriftReport returns drift analysis for all orders in a date range.
func (h *DriftHandler) GetDriftReport(c *gin.Context) {
	var req request.ListDrifts
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

	svc := h.getService(c)
	results, err := svc.DetectForPeriod(c.Request.Context(), startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	items := make([]response.DriftResult, len(results))
	for i, r := range results {
		items[i] = response.DriftResultFromDomain(r)
	}

	c.JSON(http.StatusOK, gin.H{"items": items, "total_count": len(items)})
}

package handler

import (
	"net/http"

	"github.com/diogenes/costforensics/backend/internal/adapter/http/dto/request"
	"github.com/diogenes/costforensics/backend/internal/adapter/http/dto/response"
	tenantrepo "github.com/diogenes/costforensics/backend/internal/adapter/postgres/tenant"
	pnlapp "github.com/diogenes/costforensics/backend/internal/application/pnl"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// PnLHandler handles P&L and cohort report endpoints.
type PnLHandler struct{}

// NewPnLHandler creates a new PnLHandler.
func NewPnLHandler() *PnLHandler {
	return &PnLHandler{}
}

func (h *PnLHandler) getService(c *gin.Context) *pnlapp.Service {
	db := c.MustGet("tenant_db").(*gorm.DB)
	orderRepo := tenantrepo.NewOrderRepo(db)
	sellerRepo := tenantrepo.NewSellerRepo(db)
	discountRepo := tenantrepo.NewDiscountApplicationRepo(db)
	return pnlapp.NewService(orderRepo, sellerRepo, discountRepo)
}

// GetPL generates a P&L report for the specified date range and granularity.
func (h *PnLHandler) GetPL(c *gin.Context) {
	var req request.GetPL
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

	granularity := pnlapp.ParseGranularity(req.Granularity)

	svc := h.getService(c)
	rows, err := svc.GeneratePL(c.Request.Context(), startDate, endDate, granularity, req.Currency)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	items := make([]response.PLRow, len(rows))
	for i, r := range rows {
		items[i] = response.PLRowFromDomain(r)
	}

	c.JSON(http.StatusOK, gin.H{"items": items})
}

// GetCohorts generates a retention cohort report for the specified date range.
func (h *PnLHandler) GetCohorts(c *gin.Context) {
	var req request.GetCohorts
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

	granularity := pnlapp.ParseGranularity(req.Granularity)

	svc := h.getService(c)
	rows, err := svc.GenerateCohorts(c.Request.Context(), startDate, endDate, granularity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	items := make([]response.CohortRow, len(rows))
	for i, r := range rows {
		items[i] = response.CohortRowFromDomain(r)
	}

	c.JSON(http.StatusOK, gin.H{"items": items})
}

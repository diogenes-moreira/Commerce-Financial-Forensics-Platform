package handler

import (
	"net/http"

	"github.com/diogenes/costforensics/backend/internal/adapter/http/dto/request"
	"github.com/diogenes/costforensics/backend/internal/adapter/http/dto/response"
	tenantrepo "github.com/diogenes/costforensics/backend/internal/adapter/postgres/tenant"
	costreportapp "github.com/diogenes/costforensics/backend/internal/application/costreport"
	"github.com/diogenes/costforensics/backend/internal/domain/costreport"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CostReportHandler struct{}

func NewCostReportHandler() *CostReportHandler {
	return &CostReportHandler{}
}

func (h *CostReportHandler) getService(c *gin.Context) *costreportapp.Service {
	db := c.MustGet("tenant_db").(*gorm.DB)
	reportRepo := tenantrepo.NewCostReportRepo(db)
	costRecordRepo := tenantrepo.NewCostRecordRepo(db)
	return costreportapp.NewService(reportRepo, costRecordRepo)
}

func (h *CostReportHandler) Generate(c *gin.Context) {
	var req request.GenerateCostReport
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	start, err := request.ParseDate(req.PeriodStart)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid period_start"})
		return
	}
	end, err := request.ParseDate(req.PeriodEnd)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid period_end"})
		return
	}

	currency := req.Currency
	if currency == "" {
		currency = "USD"
	}

	svc := h.getService(c)
	r, err := svc.Generate(c.Request.Context(), req.Name, costreport.ReportType(req.ReportType), start, end, currency)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, response.CostReportFromDomain(r))
}

func (h *CostReportHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	svc := h.getService(c)
	r, err := svc.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, response.CostReportFromDomain(r))
}

func (h *CostReportHandler) List(c *gin.Context) {
	svc := h.getService(c)
	reports, err := svc.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	items := make([]response.CostReport, len(reports))
	for i, r := range reports {
		items[i] = response.CostReportFromDomain(r)
	}
	c.JSON(http.StatusOK, items)
}

package handler

import (
	"net/http"

	"github.com/diogenes/costforensics/backend/internal/adapter/http/dto/request"
	"github.com/diogenes/costforensics/backend/internal/adapter/http/dto/response"
	tenantrepo "github.com/diogenes/costforensics/backend/internal/adapter/postgres/tenant"
	costrecordapp "github.com/diogenes/costforensics/backend/internal/application/costrecord"
	"github.com/diogenes/costforensics/backend/internal/domain/costrecord"
	"github.com/diogenes/costforensics/backend/internal/domain/shared"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CostRecordHandler struct{}

func NewCostRecordHandler() *CostRecordHandler {
	return &CostRecordHandler{}
}

func (h *CostRecordHandler) getService(c *gin.Context) *costrecordapp.Service {
	db := c.MustGet("tenant_db").(*gorm.DB)
	repo := tenantrepo.NewCostRecordRepo(db)
	return costrecordapp.NewService(repo)
}

func (h *CostRecordHandler) Create(c *gin.Context) {
	var req request.CreateCostRecord
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	accountID, err := uuid.Parse(req.CloudAccountID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid cloud_account_id"})
		return
	}

	currency := req.Currency
	if currency == "" {
		currency = "USD"
	}
	amount, err := shared.NewMoney(req.AmountCents, currency)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	usageDate, err := request.ParseDate(req.UsageDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid usage_date, expected YYYY-MM-DD"})
		return
	}

	svc := h.getService(c)
	r, err := svc.Create(c.Request.Context(), accountID, req.Service, amount, usageDate)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, response.CostRecordFromDomain(r))
}

func (h *CostRecordHandler) GetByID(c *gin.Context) {
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
	c.JSON(http.StatusOK, response.CostRecordFromDomain(r))
}

func (h *CostRecordHandler) List(c *gin.Context) {
	var req request.ListCostRecords
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	filter := costrecord.ListFilter{
		Service:    req.Service,
		Category:   costrecord.Category(req.Category),
		Pagination: shared.NewPagination(req.Page, req.PageSize),
	}

	if req.CloudAccountID != "" {
		id, err := uuid.Parse(req.CloudAccountID)
		if err == nil {
			filter.CloudAccountID = &id
		}
	}
	if req.StartDate != "" {
		if t, err := request.ParseDate(req.StartDate); err == nil {
			filter.StartDate = &t
		}
	}
	if req.EndDate != "" {
		if t, err := request.ParseDate(req.EndDate); err == nil {
			filter.EndDate = &t
		}
	}

	svc := h.getService(c)
	result, err := svc.List(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	items := make([]response.CostRecord, len(result.Items))
	for i, r := range result.Items {
		items[i] = response.CostRecordFromDomain(r)
	}

	c.JSON(http.StatusOK, response.Paginated[response.CostRecord]{
		Items:      items,
		TotalCount: result.TotalCount,
		Page:       result.Page,
		PageSize:   result.PageSize,
	})
}

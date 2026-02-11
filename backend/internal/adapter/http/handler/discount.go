package handler

import (
	"net/http"
	"time"

	"github.com/diogenes/costforensics/backend/internal/adapter/http/dto/request"
	"github.com/diogenes/costforensics/backend/internal/adapter/http/dto/response"
	tenantrepo "github.com/diogenes/costforensics/backend/internal/adapter/postgres/tenant"
	discountapp "github.com/diogenes/costforensics/backend/internal/application/discount"
	"github.com/diogenes/costforensics/backend/internal/domain/discount"
	"github.com/diogenes/costforensics/backend/internal/domain/shared"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DiscountHandler struct{}

func NewDiscountHandler() *DiscountHandler {
	return &DiscountHandler{}
}

func (h *DiscountHandler) getService(c *gin.Context) *discountapp.Service {
	db := c.MustGet("tenant_db").(*gorm.DB)
	ruleRepo := tenantrepo.NewPromotionRuleRepo(db)
	appRepo := tenantrepo.NewDiscountApplicationRepo(db)
	return discountapp.NewService(ruleRepo, appRepo)
}

func (h *DiscountHandler) CreateRule(c *gin.Context) {
	var req request.CreatePromotionRule
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	validFrom, err := time.Parse(time.RFC3339, req.ValidFrom)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid valid_from format, use RFC3339"})
		return
	}
	validTo, err := time.Parse(time.RFC3339, req.ValidTo)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid valid_to format, use RFC3339"})
		return
	}

	svc := h.getService(c)
	rule, err := svc.CreateRule(
		c.Request.Context(),
		req.Name, req.RuleType, req.Value,
		req.Conditions, req.FundingSource, req.FundingPct,
		req.MaxUsageCount, validFrom, validTo,
	)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, response.PromotionRuleFromDomain(rule))
}

func (h *DiscountHandler) ListRules(c *gin.Context) {
	var req request.ListPromotionRules
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	filter := discount.RuleListFilter{
		Status:     discount.RuleStatus(req.Status),
		Pagination: shared.NewPagination(req.Page, req.PageSize),
	}

	svc := h.getService(c)
	result, err := svc.ListRules(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	items := make([]response.PromotionRule, len(result.Items))
	for i, r := range result.Items {
		items[i] = response.PromotionRuleFromDomain(r)
	}
	c.JSON(http.StatusOK, response.Paginated[response.PromotionRule]{
		Items: items, TotalCount: result.TotalCount,
		Page: result.Page, PageSize: result.PageSize,
	})
}

func (h *DiscountHandler) GetRule(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	svc := h.getService(c)
	rule, err := svc.GetRule(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, response.PromotionRuleFromDomain(rule))
}

func (h *DiscountHandler) UpdateRule(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req request.CreatePromotionRule
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_ = id
	// For update, we get the existing rule and modify it
	// Since PromotionRule has no general Update method in the domain,
	// we disable the old one and create a new version
	c.JSON(http.StatusNotImplemented, gin.H{"error": "use disable + create for rule updates"})
}

func (h *DiscountHandler) DisableRule(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	svc := h.getService(c)
	rule, err := svc.DisableRule(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, response.PromotionRuleFromDomain(rule))
}

func (h *DiscountHandler) ListApplications(c *gin.Context) {
	orderIDStr := c.Query("order_id")
	if orderIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "order_id is required"})
		return
	}
	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order_id"})
		return
	}

	svc := h.getService(c)
	apps, err := svc.ListApplicationsByOrder(c.Request.Context(), orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	items := make([]response.DiscountApplication, len(apps))
	for i, a := range apps {
		items[i] = response.DiscountApplicationFromDomain(a)
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

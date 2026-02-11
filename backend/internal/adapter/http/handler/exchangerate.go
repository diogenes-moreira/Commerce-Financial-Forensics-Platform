package handler

import (
	"net/http"

	"github.com/diogenes/costforensics/backend/internal/adapter/http/dto/request"
	"github.com/diogenes/costforensics/backend/internal/adapter/http/dto/response"
	tenantrepo "github.com/diogenes/costforensics/backend/internal/adapter/postgres/tenant"
	exchangerateapp "github.com/diogenes/costforensics/backend/internal/application/exchangerate"
	"github.com/diogenes/costforensics/backend/internal/domain/exchangerate"
	"github.com/diogenes/costforensics/backend/internal/domain/shared"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ExchangeRateHandler struct{}

func NewExchangeRateHandler() *ExchangeRateHandler {
	return &ExchangeRateHandler{}
}

func (h *ExchangeRateHandler) getService(c *gin.Context) *exchangerateapp.Service {
	db := c.MustGet("tenant_db").(*gorm.DB)
	repo := tenantrepo.NewExchangeRateRepo(db)
	return exchangerateapp.NewService(repo)
}

func (h *ExchangeRateHandler) Create(c *gin.Context) {
	var req request.CreateExchangeRate
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	effectiveDate, err := request.ParseDate(req.EffectiveDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid effective_date"})
		return
	}

	source := req.Source
	if source == "" {
		source = "manual"
	}

	svc := h.getService(c)
	rate, err := svc.Create(c.Request.Context(), req.BaseCurrency, req.QuoteCurrency, req.Rate, source, effectiveDate)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, response.ExchangeRateFromDomain(rate))
}

func (h *ExchangeRateHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	svc := h.getService(c)
	rate, err := svc.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, response.ExchangeRateFromDomain(rate))
}

func (h *ExchangeRateHandler) GetLatest(c *gin.Context) {
	baseCurrency := c.Query("base_currency")
	quoteCurrency := c.Query("quote_currency")
	if baseCurrency == "" || quoteCurrency == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "base_currency and quote_currency are required"})
		return
	}

	svc := h.getService(c)
	rate, err := svc.GetLatest(c.Request.Context(), baseCurrency, quoteCurrency)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, response.ExchangeRateFromDomain(rate))
}

func (h *ExchangeRateHandler) List(c *gin.Context) {
	var req request.ListExchangeRates
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	filter := exchangerate.ListFilter{
		BaseCurrency:  req.BaseCurrency,
		QuoteCurrency: req.QuoteCurrency,
		Pagination:    shared.NewPagination(req.Page, req.PageSize),
	}
	if req.StartDate != "" {
		t, err := request.ParseDate(req.StartDate)
		if err == nil {
			filter.StartDate = &t
		}
	}
	if req.EndDate != "" {
		t, err := request.ParseDate(req.EndDate)
		if err == nil {
			filter.EndDate = &t
		}
	}

	svc := h.getService(c)
	result, err := svc.List(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	items := make([]response.ExchangeRate, len(result.Items))
	for i, r := range result.Items {
		items[i] = response.ExchangeRateFromDomain(r)
	}
	c.JSON(http.StatusOK, response.Paginated[response.ExchangeRate]{
		Items: items, TotalCount: result.TotalCount,
		Page: result.Page, PageSize: result.PageSize,
	})
}

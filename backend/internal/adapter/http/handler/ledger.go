package handler

import (
	"net/http"
	"time"

	"github.com/diogenes/costforensics/backend/internal/adapter/http/dto/request"
	"github.com/diogenes/costforensics/backend/internal/adapter/http/dto/response"
	tenantrepo "github.com/diogenes/costforensics/backend/internal/adapter/postgres/tenant"
	ledgerapp "github.com/diogenes/costforensics/backend/internal/application/ledger"
	"github.com/diogenes/costforensics/backend/internal/domain/ledger"
	"github.com/diogenes/costforensics/backend/internal/domain/shared"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LedgerHandler struct{}

func NewLedgerHandler() *LedgerHandler {
	return &LedgerHandler{}
}

func (h *LedgerHandler) getService(c *gin.Context) *ledgerapp.Service {
	db := c.MustGet("tenant_db").(*gorm.DB)
	repo := tenantrepo.NewLedgerRepo(db)
	return ledgerapp.NewService(repo)
}

func (h *LedgerHandler) ListEntries(c *gin.Context) {
	var req request.ListLedgerEntries
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	svc := h.getService(c)

	// If order_id is provided, list by order
	if req.OrderID != "" {
		orderID, err := uuid.Parse(req.OrderID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order_id"})
			return
		}
		entries, err := svc.ListByOrder(c.Request.Context(), orderID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		items := make([]response.LedgerEntry, len(entries))
		for i, e := range entries {
			items[i] = response.LedgerEntryFromDomain(e)
		}
		c.JSON(http.StatusOK, gin.H{"items": items})
		return
	}

	// Otherwise list by account code
	if req.AccountCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "order_id or account_code is required"})
		return
	}

	var dateRange ledger.DateRange
	if req.StartDate != "" {
		t, err := time.Parse("2006-01-02", req.StartDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start_date format, use YYYY-MM-DD"})
			return
		}
		dateRange.Start = t
	}
	if req.EndDate != "" {
		t, err := time.Parse("2006-01-02", req.EndDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end_date format, use YYYY-MM-DD"})
			return
		}
		dateRange.End = t
	}

	pagination := shared.NewPagination(req.Page, req.PageSize)
	result, err := svc.ListByAccount(c.Request.Context(), req.AccountCode, dateRange, pagination)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	items := make([]response.LedgerEntry, len(result.Items))
	for i, e := range result.Items {
		items[i] = response.LedgerEntryFromDomain(e)
	}
	c.JSON(http.StatusOK, response.Paginated[response.LedgerEntry]{
		Items: items, TotalCount: result.TotalCount,
		Page: result.Page, PageSize: result.PageSize,
	})
}

func (h *LedgerHandler) GetBalance(c *gin.Context) {
	var req request.GetBalance
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var dateRange ledger.DateRange
	startDateStr := ""
	endDateStr := ""
	if req.StartDate != "" {
		t, err := time.Parse("2006-01-02", req.StartDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start_date format, use YYYY-MM-DD"})
			return
		}
		dateRange.Start = t
		startDateStr = req.StartDate
	}
	if req.EndDate != "" {
		t, err := time.Parse("2006-01-02", req.EndDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end_date format, use YYYY-MM-DD"})
			return
		}
		dateRange.End = t
		endDateStr = req.EndDate
	}

	svc := h.getService(c)
	balance, err := svc.BalanceByAccount(c.Request.Context(), req.AccountCode, dateRange)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response.BalanceResponse{
		AccountCode:  req.AccountCode,
		BalanceCents: balance,
		StartDate:    startDateStr,
		EndDate:      endDateStr,
	})
}

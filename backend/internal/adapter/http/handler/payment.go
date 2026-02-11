package handler

import (
	"net/http"

	"github.com/diogenes/costforensics/backend/internal/adapter/http/dto/request"
	"github.com/diogenes/costforensics/backend/internal/adapter/http/dto/response"
	tenantrepo "github.com/diogenes/costforensics/backend/internal/adapter/postgres/tenant"
	paymentapp "github.com/diogenes/costforensics/backend/internal/application/payment"
	"github.com/diogenes/costforensics/backend/internal/domain/payment"
	"github.com/diogenes/costforensics/backend/internal/domain/shared"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PaymentHandler struct{}

func NewPaymentHandler() *PaymentHandler {
	return &PaymentHandler{}
}

func (h *PaymentHandler) getService(c *gin.Context) *paymentapp.Service {
	db := c.MustGet("tenant_db").(*gorm.DB)
	repo := tenantrepo.NewPaymentRepo(db)
	return paymentapp.NewService(repo)
}

func (h *PaymentHandler) Create(c *gin.Context) {
	var req request.CreatePayment
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	orderID, err := uuid.Parse(req.OrderID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order_id"})
		return
	}
	counterpartyID, err := uuid.Parse(req.CounterpartyID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid counterparty_id"})
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

	svc := h.getService(c)
	p, err := svc.Create(
		c.Request.Context(),
		orderID, payment.Direction(req.Direction), counterpartyID,
		amount, req.Method, req.ExternalRef,
	)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, response.PaymentFromDomain(p))
}

func (h *PaymentHandler) GetByID(c *gin.Context) {
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
	c.JSON(http.StatusOK, response.PaymentFromDomain(p))
}

func (h *PaymentHandler) List(c *gin.Context) {
	var req request.ListPayments
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	filter := payment.ListFilter{
		Direction:  payment.Direction(req.Direction),
		Status:     payment.Status(req.Status),
		Pagination: shared.NewPagination(req.Page, req.PageSize),
	}
	if req.OrderID != "" {
		id, err := uuid.Parse(req.OrderID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order_id"})
			return
		}
		filter.OrderID = &id
	}
	if req.CounterpartyID != "" {
		id, err := uuid.Parse(req.CounterpartyID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid counterparty_id"})
			return
		}
		filter.CounterpartyID = &id
	}

	svc := h.getService(c)
	result, err := svc.List(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	items := make([]response.Payment, len(result.Items))
	for i, p := range result.Items {
		items[i] = response.PaymentFromDomain(p)
	}
	c.JSON(http.StatusOK, response.Paginated[response.Payment]{
		Items: items, TotalCount: result.TotalCount,
		Page: result.Page, PageSize: result.PageSize,
	})
}

func (h *PaymentHandler) ListByOrder(c *gin.Context) {
	orderID, err := uuid.Parse(c.Param("order_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order_id"})
		return
	}
	svc := h.getService(c)
	payments, err := svc.ListByOrder(c.Request.Context(), orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	items := make([]response.Payment, len(payments))
	for i, p := range payments {
		items[i] = response.PaymentFromDomain(p)
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h *PaymentHandler) MarkProcessed(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	svc := h.getService(c)
	p, err := svc.MarkProcessed(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, response.PaymentFromDomain(p))
}

func (h *PaymentHandler) MarkFailed(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req request.FailPayment
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	svc := h.getService(c)
	p, err := svc.MarkFailed(c.Request.Context(), id, req.Reason)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, response.PaymentFromDomain(p))
}

func (h *PaymentHandler) MarkRefunded(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	svc := h.getService(c)
	p, err := svc.MarkRefunded(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, response.PaymentFromDomain(p))
}

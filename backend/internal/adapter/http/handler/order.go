package handler

import (
	"net/http"

	"github.com/diogenes/costforensics/backend/internal/adapter/http/dto/request"
	"github.com/diogenes/costforensics/backend/internal/adapter/http/dto/response"
	tenantrepo "github.com/diogenes/costforensics/backend/internal/adapter/postgres/tenant"
	orderapp "github.com/diogenes/costforensics/backend/internal/application/order"
	"github.com/diogenes/costforensics/backend/internal/domain/order"
	"github.com/diogenes/costforensics/backend/internal/domain/shared"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderHandler struct{}

func NewOrderHandler() *OrderHandler {
	return &OrderHandler{}
}

func (h *OrderHandler) getService(c *gin.Context) *orderapp.Service {
	db := c.MustGet("tenant_db").(*gorm.DB)
	repo := tenantrepo.NewOrderRepo(db)
	return orderapp.NewService(repo)
}

func (h *OrderHandler) Create(c *gin.Context) {
	var req request.CreateOrder
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sellerID, err := uuid.Parse(req.SellerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid seller_id"})
		return
	}
	customerID, err := uuid.Parse(req.CustomerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid customer_id"})
		return
	}

	svc := h.getService(c)
	o, err := svc.Create(c.Request.Context(), req.ExternalID, sellerID, customerID, req.Currency)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, response.OrderFromDomain(o))
}

func (h *OrderHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	svc := h.getService(c)
	o, err := svc.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, response.OrderFromDomain(o))
}

func (h *OrderHandler) List(c *gin.Context) {
	var req request.ListOrders
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	filter := order.ListFilter{
		Status:     order.Status(req.Status),
		Pagination: shared.NewPagination(req.Page, req.PageSize),
	}
	if req.SellerID != "" {
		id, err := uuid.Parse(req.SellerID)
		if err == nil {
			filter.SellerID = &id
		}
	}
	if req.CustomerID != "" {
		id, err := uuid.Parse(req.CustomerID)
		if err == nil {
			filter.CustomerID = &id
		}
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

	items := make([]response.Order, len(result.Items))
	for i, o := range result.Items {
		items[i] = response.OrderFromDomain(o)
	}
	c.JSON(http.StatusOK, response.Paginated[response.Order]{
		Items: items, TotalCount: result.TotalCount,
		Page: result.Page, PageSize: result.PageSize,
	})
}

func (h *OrderHandler) Update(c *gin.Context) {
	// Orders are updated through specific action endpoints (confirm, ship, etc.)
	c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "use action endpoints to update orders"})
}

func (h *OrderHandler) AddItem(c *gin.Context) {
	orderID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order id"})
		return
	}
	var req request.AddOrderItem
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	productID, err := uuid.Parse(req.ProductID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product_id"})
		return
	}

	currency := req.Currency
	if currency == "" {
		currency = "USD"
	}
	unitPrice, err := shared.NewMoney(req.UnitPriceCents, currency)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	unitCost, err := shared.NewMoney(req.UnitCostCents, currency)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	svc := h.getService(c)
	o, err := svc.AddItem(c.Request.Context(), orderID, productID, req.SKU, req.ProductName, req.Quantity, unitPrice, unitCost)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, response.OrderFromDomain(o))
}

func (h *OrderHandler) Confirm(c *gin.Context) {
	h.statusTransition(c, func(svc *orderapp.Service, ctx *gin.Context, id uuid.UUID) (*order.Order, error) {
		return svc.Confirm(ctx.Request.Context(), id)
	})
}

func (h *OrderHandler) Ship(c *gin.Context) {
	h.statusTransition(c, func(svc *orderapp.Service, ctx *gin.Context, id uuid.UUID) (*order.Order, error) {
		return svc.Ship(ctx.Request.Context(), id)
	})
}

func (h *OrderHandler) Deliver(c *gin.Context) {
	h.statusTransition(c, func(svc *orderapp.Service, ctx *gin.Context, id uuid.UUID) (*order.Order, error) {
		return svc.Deliver(ctx.Request.Context(), id)
	})
}

func (h *OrderHandler) Cancel(c *gin.Context) {
	h.statusTransition(c, func(svc *orderapp.Service, ctx *gin.Context, id uuid.UUID) (*order.Order, error) {
		return svc.Cancel(ctx.Request.Context(), id)
	})
}

func (h *OrderHandler) Refund(c *gin.Context) {
	h.statusTransition(c, func(svc *orderapp.Service, ctx *gin.Context, id uuid.UUID) (*order.Order, error) {
		return svc.Refund(ctx.Request.Context(), id)
	})
}

func (h *OrderHandler) statusTransition(c *gin.Context, fn func(*orderapp.Service, *gin.Context, uuid.UUID) (*order.Order, error)) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	svc := h.getService(c)
	o, err := fn(svc, c, id)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, response.OrderFromDomain(o))
}

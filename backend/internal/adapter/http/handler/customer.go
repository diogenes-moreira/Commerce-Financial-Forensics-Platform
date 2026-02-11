package handler

import (
	"net/http"

	"github.com/diogenes/costforensics/backend/internal/adapter/http/dto/request"
	"github.com/diogenes/costforensics/backend/internal/adapter/http/dto/response"
	tenantrepo "github.com/diogenes/costforensics/backend/internal/adapter/postgres/tenant"
	customerapp "github.com/diogenes/costforensics/backend/internal/application/customer"
	"github.com/diogenes/costforensics/backend/internal/domain/customer"
	"github.com/diogenes/costforensics/backend/internal/domain/shared"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CustomerHandler struct{}

func NewCustomerHandler() *CustomerHandler {
	return &CustomerHandler{}
}

func (h *CustomerHandler) getService(c *gin.Context) *customerapp.Service {
	db := c.MustGet("tenant_db").(*gorm.DB)
	repo := tenantrepo.NewCustomerRepo(db)
	return customerapp.NewService(repo)
}

func (h *CustomerHandler) Create(c *gin.Context) {
	var req request.CreateCustomer
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	svc := h.getService(c)
	cust, err := svc.Create(c.Request.Context(), req.ExternalID, req.Name, req.Email, req.Segment)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, response.CustomerFromDomain(cust))
}

func (h *CustomerHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	svc := h.getService(c)
	cust, err := svc.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, response.CustomerFromDomain(cust))
}

func (h *CustomerHandler) List(c *gin.Context) {
	var req request.ListCustomers
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	filter := customer.ListFilter{
		Segment:    req.Segment,
		Pagination: shared.NewPagination(req.Page, req.PageSize),
	}

	svc := h.getService(c)
	result, err := svc.List(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	items := make([]response.Customer, len(result.Items))
	for i, cust := range result.Items {
		items[i] = response.CustomerFromDomain(cust)
	}
	c.JSON(http.StatusOK, response.Paginated[response.Customer]{
		Items: items, TotalCount: result.TotalCount,
		Page: result.Page, PageSize: result.PageSize,
	})
}

func (h *CustomerHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req request.UpdateCustomer
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	svc := h.getService(c)
	cust, err := svc.Update(c.Request.Context(), id, req.Name, req.Email, req.Segment)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, response.CustomerFromDomain(cust))
}

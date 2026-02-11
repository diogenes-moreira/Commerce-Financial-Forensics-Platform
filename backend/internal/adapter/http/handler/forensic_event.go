package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/diogenes/costforensics/backend/internal/adapter/http/dto/response"
	tenantrepo "github.com/diogenes/costforensics/backend/internal/adapter/postgres/tenant"
	forensicapp "github.com/diogenes/costforensics/backend/internal/application/forensicevent"
	"github.com/diogenes/costforensics/backend/internal/domain/shared"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ForensicEventHandler struct{}

func NewForensicEventHandler() *ForensicEventHandler {
	return &ForensicEventHandler{}
}

func (h *ForensicEventHandler) getService(c *gin.Context) *forensicapp.Service {
	db := c.MustGet("tenant_db").(*gorm.DB)
	repo := tenantrepo.NewForensicEventRepo(db)
	return forensicapp.NewService(repo)
}

func (h *ForensicEventHandler) ListEvents(c *gin.Context) {
	entityType := c.Query("entity_type")
	entityIDStr := c.Query("entity_id")
	eventType := c.Query("event_type")
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")
	page := queryParamInt(c, "page", 0)
	pageSize := queryParamInt(c, "page_size", 0)

	svc := h.getService(c)

	// If entity_type and entity_id are provided, list by entity
	if entityType != "" && entityIDStr != "" {
		entityID, err := uuid.Parse(entityIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid entity_id"})
			return
		}
		events, err := svc.ListByEntity(c.Request.Context(), entityType, entityID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		items := make([]response.ForensicEvent, len(events))
		for i, e := range events {
			items[i] = response.ForensicEventFromDomain(e)
		}
		c.JSON(http.StatusOK, gin.H{"items": items})
		return
	}

	var start, end time.Time
	if startDateStr != "" {
		t, err := time.Parse("2006-01-02", startDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start_date format, use YYYY-MM-DD"})
			return
		}
		start = t
	}
	if endDateStr != "" {
		t, err := time.Parse("2006-01-02", endDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end_date format, use YYYY-MM-DD"})
			return
		}
		end = t
	}

	pagination := shared.NewPagination(page, pageSize)

	// If event_type is provided, list by type
	if eventType != "" {
		result, err := svc.ListByType(c.Request.Context(), eventType, start, end, pagination)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		items := make([]response.ForensicEvent, len(result.Items))
		for i, e := range result.Items {
			items[i] = response.ForensicEventFromDomain(e)
		}
		c.JSON(http.StatusOK, response.Paginated[response.ForensicEvent]{
			Items: items, TotalCount: result.TotalCount,
			Page: result.Page, PageSize: result.PageSize,
		})
		return
	}

	// Default: list by date range
	result, err := svc.ListByDateRange(c.Request.Context(), start, end, pagination)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	items := make([]response.ForensicEvent, len(result.Items))
	for i, e := range result.Items {
		items[i] = response.ForensicEventFromDomain(e)
	}
	c.JSON(http.StatusOK, response.Paginated[response.ForensicEvent]{
		Items: items, TotalCount: result.TotalCount,
		Page: result.Page, PageSize: result.PageSize,
	})
}

func queryParamInt(c *gin.Context, key string, defaultVal int) int {
	val := c.Query(key)
	if val == "" {
		return defaultVal
	}
	i, err := strconv.Atoi(val)
	if err != nil {
		return defaultVal
	}
	return i
}

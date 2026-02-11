package middleware

import (
	"net/http"

	"github.com/diogenes/costforensics/backend/internal/domain/tenant"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func TenantResolver(repo tenant.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantIDVal, exists := c.Get("tenant_id")
		if !exists {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "tenant_id not found in context"})
			return
		}

		tenantID, ok := tenantIDVal.(uuid.UUID)
		if !ok {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid tenant_id"})
			return
		}

		t, err := repo.GetByID(c.Request.Context(), tenantID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "tenant not found"})
			return
		}

		if !t.IsActive() {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "tenant is suspended"})
			return
		}

		c.Set("tenant", t)
		c.Set("tenant_db_name", t.DBName())
		c.Next()
	}
}

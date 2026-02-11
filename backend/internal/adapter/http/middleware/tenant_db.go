package middleware

import (
	"net/http"

	"github.com/diogenes/costforensics/backend/internal/platform/database"
	"github.com/gin-gonic/gin"
)

func TenantDB(dbManager *database.TenantDBManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		dbNameVal, exists := c.Get("tenant_db_name")
		if !exists {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "tenant db name not resolved"})
			return
		}

		dbName, ok := dbNameVal.(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "invalid tenant db name"})
			return
		}

		db, err := dbManager.GetDB(dbName)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to connect to tenant database"})
			return
		}

		c.Set("tenant_db", db)
		c.Next()
	}
}

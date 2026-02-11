package http

import (
	"github.com/diogenes/costforensics/backend/internal/adapter/http/handler"
	"github.com/diogenes/costforensics/backend/internal/adapter/http/middleware"
	"github.com/diogenes/costforensics/backend/internal/domain/tenant"
	"github.com/diogenes/costforensics/backend/internal/platform/auth"
	"github.com/diogenes/costforensics/backend/internal/platform/database"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type RouterDeps struct {
	Log          *logrus.Logger
	SystemDB     *gorm.DB
	DBManager    *database.TenantDBManager
	JWTValidator *auth.JWTValidator
	TenantRepo   tenant.Repository

	// Handlers
	HealthHandler       *handler.HealthHandler
	TenantHandler       *handler.TenantHandler
	CloudAccountHandler *handler.CloudAccountHandler
	CostRecordHandler   *handler.CostRecordHandler
	BudgetHandler       *handler.BudgetHandler
	AnomalyHandler      *handler.AnomalyHandler
	CostReportHandler   *handler.CostReportHandler
}

func NewRouter(deps RouterDeps) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	// Global middleware chain: Recovery → RequestID → Logging → CORS
	r.Use(
		middleware.Recovery(deps.Log),
		middleware.RequestID(),
		middleware.Logging(deps.Log),
		middleware.CORS(),
	)

	// Public routes
	r.GET("/healthz", deps.HealthHandler.Healthz)

	// Admin routes (for tenant management — no tenant scoping needed)
	admin := r.Group("/api/v1/admin")
	admin.Use(middleware.Auth(deps.JWTValidator))
	{
		admin.POST("/tenants", deps.TenantHandler.Create)
		admin.GET("/tenants", deps.TenantHandler.List)
		admin.GET("/tenants/:id", deps.TenantHandler.GetByID)
		admin.PUT("/tenants/:id", deps.TenantHandler.Update)
		admin.DELETE("/tenants/:id", deps.TenantHandler.Delete)
	}

	// Tenant-scoped routes: Auth → TenantResolver → TenantDB → Handler
	api := r.Group("/api/v1")
	api.Use(
		middleware.Auth(deps.JWTValidator),
		middleware.TenantResolver(deps.TenantRepo),
		middleware.TenantDB(deps.DBManager),
	)
	{
		// Cloud Accounts
		api.POST("/cloud-accounts", deps.CloudAccountHandler.Create)
		api.GET("/cloud-accounts", deps.CloudAccountHandler.List)
		api.GET("/cloud-accounts/:id", deps.CloudAccountHandler.GetByID)
		api.PUT("/cloud-accounts/:id", deps.CloudAccountHandler.Update)
		api.DELETE("/cloud-accounts/:id", deps.CloudAccountHandler.Delete)
		api.POST("/cloud-accounts/:id/sync", deps.CloudAccountHandler.BeginSync)

		// Cost Records
		api.POST("/cost-records", deps.CostRecordHandler.Create)
		api.GET("/cost-records", deps.CostRecordHandler.List)
		api.GET("/cost-records/:id", deps.CostRecordHandler.GetByID)

		// Budgets
		api.POST("/budgets", deps.BudgetHandler.Create)
		api.GET("/budgets", deps.BudgetHandler.List)
		api.GET("/budgets/:id", deps.BudgetHandler.GetByID)
		api.PUT("/budgets/:id", deps.BudgetHandler.Update)
		api.POST("/budgets/:id/spend", deps.BudgetHandler.RecordSpend)

		// Anomalies
		api.GET("/anomalies", deps.AnomalyHandler.List)
		api.GET("/anomalies/:id", deps.AnomalyHandler.GetByID)
		api.POST("/anomalies/:id/resolve", deps.AnomalyHandler.Resolve)

		// Cost Reports
		api.POST("/cost-reports", deps.CostReportHandler.Generate)
		api.GET("/cost-reports", deps.CostReportHandler.List)
		api.GET("/cost-reports/:id", deps.CostReportHandler.GetByID)
	}

	return r
}

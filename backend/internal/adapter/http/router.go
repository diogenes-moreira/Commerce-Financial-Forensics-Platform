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

	// Commerce handlers
	ProductHandler     *handler.ProductHandler
	SellerHandler      *handler.SellerHandler
	CustomerHandler    *handler.CustomerHandler
	OrderHandler       *handler.OrderHandler

	// Financial handlers
	LedgerHandler       *handler.LedgerHandler
	ForensicEventHandler *handler.ForensicEventHandler
	DiscountHandler     *handler.DiscountHandler
	PaymentHandler      *handler.PaymentHandler
	ExchangeRateHandler *handler.ExchangeRateHandler

	// Analytics handlers
	MarginHandler *handler.MarginHandler
	DriftHandler  *handler.DriftHandler
	PnLHandler    *handler.PnLHandler

	// Import handler
	ImportHandler *handler.ImportHandler

	// Integration handler
	IntegrationHandler *handler.IntegrationHandler
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

		// Products
		api.POST("/products", deps.ProductHandler.Create)
		api.GET("/products", deps.ProductHandler.List)
		api.GET("/products/:id", deps.ProductHandler.GetByID)
		api.PUT("/products/:id", deps.ProductHandler.Update)

		// Sellers
		api.POST("/sellers", deps.SellerHandler.Create)
		api.GET("/sellers", deps.SellerHandler.List)
		api.GET("/sellers/:id", deps.SellerHandler.GetByID)
		api.PUT("/sellers/:id", deps.SellerHandler.Update)

		// Customers
		api.POST("/customers", deps.CustomerHandler.Create)
		api.GET("/customers", deps.CustomerHandler.List)
		api.GET("/customers/:id", deps.CustomerHandler.GetByID)
		api.PUT("/customers/:id", deps.CustomerHandler.Update)

		// Orders
		api.POST("/orders", deps.OrderHandler.Create)
		api.GET("/orders", deps.OrderHandler.List)
		api.GET("/orders/:id", deps.OrderHandler.GetByID)
		api.POST("/orders/:id/items", deps.OrderHandler.AddItem)
		api.POST("/orders/:id/confirm", deps.OrderHandler.Confirm)
		api.POST("/orders/:id/ship", deps.OrderHandler.Ship)
		api.POST("/orders/:id/deliver", deps.OrderHandler.Deliver)
		api.POST("/orders/:id/cancel", deps.OrderHandler.Cancel)
		api.POST("/orders/:id/refund", deps.OrderHandler.Refund)

		// Ledger
		api.GET("/ledger/entries", deps.LedgerHandler.ListEntries)
		api.GET("/ledger/balance", deps.LedgerHandler.GetBalance)

		// Forensic Events
		api.GET("/events", deps.ForensicEventHandler.ListEvents)

		// Promotion Rules & Discounts
		api.POST("/promotion-rules", deps.DiscountHandler.CreateRule)
		api.GET("/promotion-rules", deps.DiscountHandler.ListRules)
		api.GET("/promotion-rules/:id", deps.DiscountHandler.GetRule)
		api.POST("/promotion-rules/:id/disable", deps.DiscountHandler.DisableRule)
		api.GET("/discount-applications", deps.DiscountHandler.ListApplications)

		// Payments
		api.POST("/payments", deps.PaymentHandler.Create)
		api.GET("/payments", deps.PaymentHandler.List)
		api.GET("/payments/:id", deps.PaymentHandler.GetByID)
		api.POST("/payments/:id/process", deps.PaymentHandler.MarkProcessed)
		api.POST("/payments/:id/fail", deps.PaymentHandler.MarkFailed)
		api.POST("/payments/:id/refund", deps.PaymentHandler.MarkRefunded)

		// Exchange Rates
		api.POST("/exchange-rates", deps.ExchangeRateHandler.Create)
		api.GET("/exchange-rates", deps.ExchangeRateHandler.List)
		api.GET("/exchange-rates/latest", deps.ExchangeRateHandler.GetLatest)
		api.GET("/exchange-rates/:id", deps.ExchangeRateHandler.GetByID)

		// Margins
		api.GET("/orders/:id/margin", deps.MarginHandler.GetOrderMargin)
		api.GET("/margins", deps.MarginHandler.ListMargins)

		// Drift
		api.GET("/orders/:id/drift", deps.DriftHandler.GetOrderDrift)
		api.GET("/drift-report", deps.DriftHandler.GetDriftReport)

		// P&L and Cohorts
		api.GET("/pnl", deps.PnLHandler.GetPL)
		api.GET("/cohorts", deps.PnLHandler.GetCohorts)

		// Imports
		api.POST("/imports/upload", deps.ImportHandler.Upload)
		api.GET("/imports", deps.ImportHandler.List)
		api.GET("/imports/:id", deps.ImportHandler.GetByID)

		// Integrations
		api.POST("/integrations", deps.IntegrationHandler.Create)
		api.GET("/integrations", deps.IntegrationHandler.List)
		api.GET("/integrations/:id", deps.IntegrationHandler.GetByID)
		api.PUT("/integrations/:id", deps.IntegrationHandler.Update)
		api.DELETE("/integrations/:id", deps.IntegrationHandler.Delete)
		api.POST("/integrations/:id/sync", deps.IntegrationHandler.TriggerSync)
		api.POST("/webhooks/:integration_id", deps.IntegrationHandler.HandleWebhook)
	}

	return r
}

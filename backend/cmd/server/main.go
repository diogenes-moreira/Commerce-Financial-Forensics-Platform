package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/diogenes/costforensics/backend/config"
	httpAdapter "github.com/diogenes/costforensics/backend/internal/adapter/http"
	"github.com/diogenes/costforensics/backend/internal/adapter/http/handler"
	"github.com/diogenes/costforensics/backend/internal/adapter/postgres/system"
	tenantapp "github.com/diogenes/costforensics/backend/internal/application/tenant"
	"github.com/diogenes/costforensics/backend/internal/platform/auth"
	"github.com/diogenes/costforensics/backend/internal/platform/database"
	"github.com/diogenes/costforensics/backend/internal/platform/logger"
)

func main() {
	// Load config
	cfg, err := config.Load("config/config.yaml")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Setup logger
	log := logger.Setup(cfg.Log)
	log.Info("Starting costForensics server")

	// Connect to system database
	systemDB, err := database.NewSystemDB(cfg.Database, log)
	if err != nil {
		log.WithError(err).Fatal("Failed to connect to system database")
	}

	// Run system migrations
	migrator := database.NewMigrator(log)
	if err := migrator.MigrateSystem(systemDB); err != nil {
		log.WithError(err).Fatal("Failed to run system migrations")
	}

	// Initialize tenant DB manager
	dbManager := database.NewTenantDBManager(cfg.Database, log)
	defer dbManager.Close()

	// Provisioner
	provisioner := database.NewProvisioner(systemDB, dbManager, migrator, log)

	// JWT validator
	jwtValidator := auth.NewJWTValidator(cfg.Auth.JWTSecret)

	// Repositories
	tenantRepo := system.NewTenantRepo(systemDB)

	// Application services
	tenantService := tenantapp.NewService(tenantRepo, provisioner, log)

	// Handlers
	healthHandler := handler.NewHealthHandler(systemDB)
	tenantHandler := handler.NewTenantHandler(tenantService)
	cloudAccountHandler := handler.NewCloudAccountHandler()
	costRecordHandler := handler.NewCostRecordHandler()
	budgetHandler := handler.NewBudgetHandler()
	anomalyHandler := handler.NewAnomalyHandler()
	costReportHandler := handler.NewCostReportHandler()

	// Router
	router := httpAdapter.NewRouter(httpAdapter.RouterDeps{
		Log:                 log,
		SystemDB:            systemDB,
		DBManager:           dbManager,
		JWTValidator:        jwtValidator,
		TenantRepo:          tenantRepo,
		HealthHandler:       healthHandler,
		TenantHandler:       tenantHandler,
		CloudAccountHandler: cloudAccountHandler,
		CostRecordHandler:   costRecordHandler,
		BudgetHandler:       budgetHandler,
		AnomalyHandler:      anomalyHandler,
		CostReportHandler:   costReportHandler,
	})

	// HTTP server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown
	go func() {
		log.WithField("port", cfg.Server.Port).Info("Server listening")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.WithError(err).Fatal("Server failed")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.WithError(err).Fatal("Server forced to shutdown")
	}

	log.Info("Server stopped")
}

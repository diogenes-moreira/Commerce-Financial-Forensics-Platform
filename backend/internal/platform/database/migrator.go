package database

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// Migrator runs schema migrations.
type Migrator struct {
	log *logrus.Logger
}

func NewMigrator(log *logrus.Logger) *Migrator {
	return &Migrator{log: log}
}

// MigrateSystem runs system database migrations (tenant registry).
func (m *Migrator) MigrateSystem(db *gorm.DB) error {
	m.log.Info("Running system database migrations")

	sql := `
	CREATE TABLE IF NOT EXISTS tenants (
		id UUID PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		slug VARCHAR(255) NOT NULL UNIQUE,
		db_name VARCHAR(255) NOT NULL UNIQUE,
		status VARCHAR(50) NOT NULL DEFAULT 'active',
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);
	CREATE INDEX IF NOT EXISTS idx_tenants_slug ON tenants(slug);
	CREATE INDEX IF NOT EXISTS idx_tenants_status ON tenants(status);
	`

	if err := db.Exec(sql).Error; err != nil {
		return fmt.Errorf("system migration failed: %w", err)
	}

	m.log.Info("System database migrations completed")
	return nil
}

// MigrateTenant runs migrations on a specific tenant database.
func (m *Migrator) MigrateTenant(db *gorm.DB) error {
	m.log.Info("Running tenant database migrations")

	sql := `
	CREATE TABLE IF NOT EXISTS cloud_accounts (
		id UUID PRIMARY KEY,
		provider VARCHAR(50) NOT NULL,
		name VARCHAR(255) NOT NULL,
		external_id VARCHAR(255) NOT NULL,
		status VARCHAR(50) NOT NULL DEFAULT 'active',
		sync_status VARCHAR(50) NOT NULL DEFAULT 'idle',
		last_synced_at TIMESTAMPTZ,
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS cost_records (
		id UUID PRIMARY KEY,
		cloud_account_id UUID NOT NULL REFERENCES cloud_accounts(id),
		service VARCHAR(255) NOT NULL,
		category VARCHAR(255) NOT NULL DEFAULT 'other',
		amount_cents BIGINT NOT NULL,
		currency VARCHAR(10) NOT NULL DEFAULT 'USD',
		usage_date DATE NOT NULL,
		tags JSONB DEFAULT '{}',
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);
	CREATE INDEX IF NOT EXISTS idx_cost_records_account ON cost_records(cloud_account_id);
	CREATE INDEX IF NOT EXISTS idx_cost_records_date ON cost_records(usage_date);
	CREATE INDEX IF NOT EXISTS idx_cost_records_service ON cost_records(service);

	CREATE TABLE IF NOT EXISTS budgets (
		id UUID PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		amount_cents BIGINT NOT NULL,
		currency VARCHAR(10) NOT NULL DEFAULT 'USD',
		spent_cents BIGINT NOT NULL DEFAULT 0,
		period_start DATE NOT NULL,
		period_end DATE NOT NULL,
		alert_threshold_pct INTEGER NOT NULL DEFAULT 80,
		status VARCHAR(50) NOT NULL DEFAULT 'active',
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS budget_alerts (
		id UUID PRIMARY KEY,
		budget_id UUID NOT NULL REFERENCES budgets(id),
		threshold_pct INTEGER NOT NULL,
		actual_pct INTEGER NOT NULL,
		message TEXT NOT NULL,
		triggered_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS anomalies (
		id UUID PRIMARY KEY,
		cloud_account_id UUID NOT NULL REFERENCES cloud_accounts(id),
		service VARCHAR(255) NOT NULL,
		expected_cents BIGINT NOT NULL,
		actual_cents BIGINT NOT NULL,
		currency VARCHAR(10) NOT NULL DEFAULT 'USD',
		deviation_pct DOUBLE PRECISION NOT NULL,
		severity VARCHAR(50) NOT NULL,
		status VARCHAR(50) NOT NULL DEFAULT 'open',
		detected_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		resolved_at TIMESTAMPTZ
	);

	CREATE TABLE IF NOT EXISTS cost_reports (
		id UUID PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		report_type VARCHAR(50) NOT NULL,
		period_start DATE NOT NULL,
		period_end DATE NOT NULL,
		total_cents BIGINT NOT NULL DEFAULT 0,
		currency VARCHAR(10) NOT NULL DEFAULT 'USD',
		breakdown JSONB DEFAULT '{}',
		generated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);
	`

	if err := db.Exec(sql).Error; err != nil {
		return fmt.Errorf("tenant migration failed: %w", err)
	}

	m.log.Info("Tenant database migrations completed")
	return nil
}

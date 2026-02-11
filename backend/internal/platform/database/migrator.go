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

	-- Commerce: Products
	CREATE TABLE IF NOT EXISTS products (
		id UUID PRIMARY KEY,
		sku VARCHAR(255) NOT NULL UNIQUE,
		ean VARCHAR(13),
		upc VARCHAR(12),
		name VARCHAR(255) NOT NULL,
		description TEXT,
		category VARCHAR(255),
		unit_cost_cents BIGINT NOT NULL,
		unit_price_cents BIGINT NOT NULL,
		currency VARCHAR(10) NOT NULL DEFAULT 'USD',
		status VARCHAR(50) NOT NULL DEFAULT 'active',
		metadata JSONB DEFAULT '{}',
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);
	CREATE INDEX IF NOT EXISTS idx_products_sku ON products(sku);
	CREATE INDEX IF NOT EXISTS idx_products_category ON products(category);
	CREATE INDEX IF NOT EXISTS idx_products_status ON products(status);

	-- Commerce: Sellers
	CREATE TABLE IF NOT EXISTS sellers (
		id UUID PRIMARY KEY,
		external_id VARCHAR(255),
		code VARCHAR(255),
		name VARCHAR(255) NOT NULL,
		email VARCHAR(255),
		commission_pct DOUBLE PRECISION NOT NULL DEFAULT 0,
		status VARCHAR(50) NOT NULL DEFAULT 'active',
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);
	CREATE INDEX IF NOT EXISTS idx_sellers_external_id ON sellers(external_id);

	-- Commerce: Customers
	CREATE TABLE IF NOT EXISTS customers (
		id UUID PRIMARY KEY,
		external_id VARCHAR(255),
		name VARCHAR(255) NOT NULL,
		email VARCHAR(255),
		segment VARCHAR(100),
		metadata JSONB DEFAULT '{}',
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);
	CREATE INDEX IF NOT EXISTS idx_customers_external_id ON customers(external_id);
	CREATE INDEX IF NOT EXISTS idx_customers_segment ON customers(segment);

	-- Commerce: Orders
	CREATE TABLE IF NOT EXISTS orders (
		id UUID PRIMARY KEY,
		external_id VARCHAR(255) NOT NULL,
		seller_id UUID NOT NULL REFERENCES sellers(id),
		customer_id UUID NOT NULL REFERENCES customers(id),
		status VARCHAR(50) NOT NULL DEFAULT 'pending',
		subtotal_cents BIGINT NOT NULL DEFAULT 0,
		discount_cents BIGINT NOT NULL DEFAULT 0,
		shipping_cents BIGINT NOT NULL DEFAULT 0,
		tax_cents BIGINT NOT NULL DEFAULT 0,
		total_cents BIGINT NOT NULL DEFAULT 0,
		currency VARCHAR(10) NOT NULL DEFAULT 'USD',
		order_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		metadata JSONB DEFAULT '{}',
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);
	CREATE INDEX IF NOT EXISTS idx_orders_external_id ON orders(external_id);
	CREATE INDEX IF NOT EXISTS idx_orders_seller ON orders(seller_id);
	CREATE INDEX IF NOT EXISTS idx_orders_customer ON orders(customer_id);
	CREATE INDEX IF NOT EXISTS idx_orders_status ON orders(status);
	CREATE INDEX IF NOT EXISTS idx_orders_date ON orders(order_date);

	-- Commerce: Order Items
	CREATE TABLE IF NOT EXISTS order_items (
		id UUID PRIMARY KEY,
		order_id UUID NOT NULL REFERENCES orders(id),
		product_id UUID NOT NULL REFERENCES products(id),
		sku VARCHAR(255) NOT NULL,
		product_name VARCHAR(255) NOT NULL,
		quantity INTEGER NOT NULL,
		unit_price_cents BIGINT NOT NULL,
		unit_cost_cents BIGINT NOT NULL,
		discount_cents BIGINT NOT NULL DEFAULT 0,
		total_cents BIGINT NOT NULL,
		currency VARCHAR(10) NOT NULL DEFAULT 'USD',
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);
	CREATE INDEX IF NOT EXISTS idx_order_items_order ON order_items(order_id);
	CREATE INDEX IF NOT EXISTS idx_order_items_product ON order_items(product_id);

	-- Financial: Ledger Entries (double-entry)
	CREATE TABLE IF NOT EXISTS ledger_entries (
		id UUID PRIMARY KEY,
		order_id UUID,
		account_code VARCHAR(100) NOT NULL,
		side VARCHAR(10) NOT NULL,
		amount_cents BIGINT NOT NULL,
		currency VARCHAR(10) NOT NULL DEFAULT 'USD',
		description TEXT,
		reference_type VARCHAR(50),
		reference_id UUID,
		effective_date TIMESTAMPTZ NOT NULL,
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);
	CREATE INDEX IF NOT EXISTS idx_ledger_order ON ledger_entries(order_id);
	CREATE INDEX IF NOT EXISTS idx_ledger_account ON ledger_entries(account_code);
	CREATE INDEX IF NOT EXISTS idx_ledger_date ON ledger_entries(effective_date);

	-- Financial: Forensic Events (immutable audit log)
	CREATE TABLE IF NOT EXISTS forensic_events (
		id UUID PRIMARY KEY,
		entity_type VARCHAR(100) NOT NULL,
		entity_id UUID NOT NULL,
		event_type VARCHAR(100) NOT NULL,
		event_time TIMESTAMPTZ NOT NULL,
		effective_time TIMESTAMPTZ NOT NULL,
		payload JSONB DEFAULT '{}',
		actor_id VARCHAR(255),
		hash_integrity VARCHAR(64) NOT NULL,
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);
	CREATE INDEX IF NOT EXISTS idx_events_entity ON forensic_events(entity_type, entity_id);
	CREATE INDEX IF NOT EXISTS idx_events_type ON forensic_events(event_type);
	CREATE INDEX IF NOT EXISTS idx_events_time ON forensic_events(event_time);

	-- Financial: Promotion Rules
	CREATE TABLE IF NOT EXISTS promotion_rules (
		id UUID PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		rule_type VARCHAR(50) NOT NULL,
		value DOUBLE PRECISION NOT NULL,
		conditions JSONB DEFAULT '{}',
		funding_source VARCHAR(50),
		funding_pct DOUBLE PRECISION NOT NULL DEFAULT 0,
		max_usage_count INTEGER NOT NULL DEFAULT 0,
		current_usage INTEGER NOT NULL DEFAULT 0,
		valid_from TIMESTAMPTZ NOT NULL,
		valid_to TIMESTAMPTZ NOT NULL,
		status VARCHAR(50) NOT NULL DEFAULT 'active',
		version INTEGER NOT NULL DEFAULT 1,
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);
	CREATE INDEX IF NOT EXISTS idx_promo_status ON promotion_rules(status);

	-- Financial: Discount Applications
	CREATE TABLE IF NOT EXISTS discount_applications (
		id UUID PRIMARY KEY,
		order_id UUID NOT NULL REFERENCES orders(id),
		promotion_rule_id UUID NOT NULL REFERENCES promotion_rules(id),
		rule_version INTEGER NOT NULL,
		discount_cents BIGINT NOT NULL,
		currency VARCHAR(10) NOT NULL DEFAULT 'USD',
		funding_source VARCHAR(50),
		seller_share_cents BIGINT NOT NULL DEFAULT 0,
		platform_share_cents BIGINT NOT NULL DEFAULT 0,
		applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);
	CREATE INDEX IF NOT EXISTS idx_discount_app_order ON discount_applications(order_id);

	-- Financial: Payments (inbound from customer + outbound to seller)
	CREATE TABLE IF NOT EXISTS payments (
		id UUID PRIMARY KEY,
		order_id UUID NOT NULL REFERENCES orders(id),
		direction VARCHAR(20) NOT NULL,
		counterparty_id UUID NOT NULL,
		amount_cents BIGINT NOT NULL,
		currency VARCHAR(10) NOT NULL DEFAULT 'USD',
		method VARCHAR(100),
		external_ref VARCHAR(255),
		status VARCHAR(50) NOT NULL DEFAULT 'pending',
		processed_at TIMESTAMPTZ,
		failure_reason TEXT,
		metadata JSONB DEFAULT '{}',
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);
	CREATE INDEX IF NOT EXISTS idx_payments_order ON payments(order_id);
	CREATE INDEX IF NOT EXISTS idx_payments_direction ON payments(direction);
	CREATE INDEX IF NOT EXISTS idx_payments_counterparty ON payments(counterparty_id);
	CREATE INDEX IF NOT EXISTS idx_payments_status ON payments(status);

	-- Financial: Exchange Rates
	CREATE TABLE IF NOT EXISTS exchange_rates (
		id UUID PRIMARY KEY,
		base_currency VARCHAR(10) NOT NULL,
		quote_currency VARCHAR(10) NOT NULL,
		rate DOUBLE PRECISION NOT NULL,
		inverse_rate DOUBLE PRECISION NOT NULL,
		source VARCHAR(50),
		effective_date DATE NOT NULL,
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);
	CREATE INDEX IF NOT EXISTS idx_fx_pair ON exchange_rates(base_currency, quote_currency);
	CREATE INDEX IF NOT EXISTS idx_fx_date ON exchange_rates(effective_date);

	-- Import: Import Jobs
	CREATE TABLE IF NOT EXISTS import_jobs (
		id UUID PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		source VARCHAR(50) NOT NULL,
		entity_type VARCHAR(100) NOT NULL,
		status VARCHAR(50) NOT NULL DEFAULT 'pending',
		total_rows INTEGER NOT NULL DEFAULT 0,
		processed_rows INTEGER NOT NULL DEFAULT 0,
		failed_rows INTEGER NOT NULL DEFAULT 0,
		error_log TEXT DEFAULT '[]',
		source_uri TEXT NOT NULL,
		started_at TIMESTAMPTZ,
		completed_at TIMESTAMPTZ,
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);
	CREATE INDEX IF NOT EXISTS idx_import_jobs_status ON import_jobs(status);
	CREATE INDEX IF NOT EXISTS idx_import_jobs_entity_type ON import_jobs(entity_type);

	-- Integrations
	CREATE TABLE IF NOT EXISTS integrations (
		id UUID PRIMARY KEY,
		platform VARCHAR(50) NOT NULL,
		name VARCHAR(255) NOT NULL,
		config TEXT DEFAULT '{}',
		status VARCHAR(50) NOT NULL DEFAULT 'active',
		sync_status VARCHAR(50) NOT NULL DEFAULT 'idle',
		last_synced_at TIMESTAMPTZ,
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);
	CREATE INDEX IF NOT EXISTS idx_integrations_platform ON integrations(platform);
	CREATE INDEX IF NOT EXISTS idx_integrations_status ON integrations(status);
	`

	if err := db.Exec(sql).Error; err != nil {
		return fmt.Errorf("tenant migration failed: %w", err)
	}

	m.log.Info("Tenant database migrations completed")
	return nil
}

package database

import (
	"fmt"
	"sync"

	"github.com/diogenes/costforensics/backend/config"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// TenantDBManager manages per-tenant database connections with lazy init
// and double-checked locking.
type TenantDBManager struct {
	cfg config.DatabaseConfig
	log *logrus.Logger
	mu  sync.RWMutex
	dbs map[string]*gorm.DB
}

func NewTenantDBManager(cfg config.DatabaseConfig, log *logrus.Logger) *TenantDBManager {
	return &TenantDBManager{
		cfg: cfg,
		log: log,
		dbs: make(map[string]*gorm.DB),
	}
}

// GetDB returns a connection for the given tenant database name.
// It uses double-checked locking for thread safety with minimal contention.
func (m *TenantDBManager) GetDB(dbName string) (*gorm.DB, error) {
	// Fast path: read lock
	m.mu.RLock()
	if db, ok := m.dbs[dbName]; ok {
		m.mu.RUnlock()
		return db, nil
	}
	m.mu.RUnlock()

	// Slow path: write lock
	m.mu.Lock()
	defer m.mu.Unlock()

	// Double-check after acquiring write lock
	if db, ok := m.dbs[dbName]; ok {
		return db, nil
	}

	dsn := m.cfg.TenantDSN(dbName)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Silent),
	})
	if err != nil {
		return nil, fmt.Errorf("connecting to tenant db %s: %w", dbName, err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("getting sql.DB for tenant %s: %w", dbName, err)
	}
	sqlDB.SetMaxOpenConns(m.cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(m.cfg.MaxIdleConns)

	m.dbs[dbName] = db
	m.log.WithField("db", dbName).Info("Opened tenant database connection")
	return db, nil
}

// Close closes all tenant database connections.
func (m *TenantDBManager) Close() {
	m.mu.Lock()
	defer m.mu.Unlock()

	for name, db := range m.dbs {
		sqlDB, err := db.DB()
		if err == nil {
			sqlDB.Close()
		}
		m.log.WithField("db", name).Info("Closed tenant database connection")
	}
	m.dbs = make(map[string]*gorm.DB)
}

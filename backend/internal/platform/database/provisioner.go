package database

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// Provisioner handles creating and dropping tenant databases.
type Provisioner struct {
	systemDB  *gorm.DB
	dbManager *TenantDBManager
	migrator  *Migrator
	log       *logrus.Logger
}

func NewProvisioner(systemDB *gorm.DB, dbManager *TenantDBManager, migrator *Migrator, log *logrus.Logger) *Provisioner {
	return &Provisioner{
		systemDB:  systemDB,
		dbManager: dbManager,
		migrator:  migrator,
		log:       log,
	}
}

// ProvisionTenantDB creates a new database for a tenant and runs migrations.
func (p *Provisioner) ProvisionTenantDB(dbName string) error {
	p.log.WithField("db", dbName).Info("Provisioning tenant database")

	// Create the database using the system connection
	sql := fmt.Sprintf("CREATE DATABASE %s", dbName)
	if err := p.systemDB.Exec(sql).Error; err != nil {
		return fmt.Errorf("creating tenant database %s: %w", dbName, err)
	}

	// Connect to the new database
	tenantDB, err := p.dbManager.GetDB(dbName)
	if err != nil {
		return fmt.Errorf("connecting to new tenant database %s: %w", dbName, err)
	}

	// Run tenant migrations
	if err := p.migrator.MigrateTenant(tenantDB); err != nil {
		return fmt.Errorf("migrating tenant database %s: %w", dbName, err)
	}

	p.log.WithField("db", dbName).Info("Tenant database provisioned")
	return nil
}

package model

import (
	"time"

	"github.com/google/uuid"
)

type Tenant struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name      string    `gorm:"type:varchar(255);not null"`
	Slug      string    `gorm:"type:varchar(255);not null;uniqueIndex"`
	DBName    string    `gorm:"column:db_name;type:varchar(255);not null;uniqueIndex"`
	Status    string    `gorm:"type:varchar(50);not null;default:'active'"`
	CreatedAt time.Time `gorm:"type:timestamptz;not null"`
	UpdatedAt time.Time `gorm:"type:timestamptz;not null"`
}

func (Tenant) TableName() string { return "tenants" }

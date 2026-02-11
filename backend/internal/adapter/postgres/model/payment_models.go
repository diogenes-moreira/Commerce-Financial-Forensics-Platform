package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type Payment struct {
	ID             uuid.UUID      `gorm:"type:uuid;primaryKey"`
	OrderID        uuid.UUID      `gorm:"type:uuid;not null"`
	Direction      string         `gorm:"type:varchar(20);not null"`
	CounterpartyID uuid.UUID      `gorm:"type:uuid;not null"`
	AmountCents    int64          `gorm:"not null"`
	Currency       string         `gorm:"type:varchar(10);not null;default:'USD'"`
	Method         string         `gorm:"type:varchar(100);not null"`
	ExternalRef    string         `gorm:"type:varchar(255)"`
	Status         string         `gorm:"type:varchar(50);not null;default:'pending'"`
	ProcessedAt    *time.Time     `gorm:"type:timestamptz"`
	FailureReason  string         `gorm:"type:text"`
	Metadata       datatypes.JSON `gorm:"type:jsonb;default:'{}'"`
	CreatedAt      time.Time      `gorm:"type:timestamptz;not null"`
	UpdatedAt      time.Time      `gorm:"type:timestamptz;not null"`
}

func (Payment) TableName() string { return "payments" }

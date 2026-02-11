package model

import (
	"time"

	"github.com/google/uuid"
)

type LedgerEntry struct {
	ID            uuid.UUID  `gorm:"type:uuid;primaryKey"`
	OrderID       *uuid.UUID `gorm:"type:uuid"`
	AccountCode   string     `gorm:"type:varchar(255);not null"`
	Side          string     `gorm:"type:varchar(10);not null"`
	AmountCents   int64      `gorm:"not null"`
	Currency      string     `gorm:"type:varchar(10);not null;default:'USD'"`
	Description   string     `gorm:"type:text"`
	ReferenceType string     `gorm:"type:varchar(100);not null"`
	ReferenceID   uuid.UUID  `gorm:"type:uuid;not null"`
	EffectiveDate time.Time  `gorm:"type:timestamptz;not null"`
	CreatedAt     time.Time  `gorm:"type:timestamptz;not null"`
}

func (LedgerEntry) TableName() string { return "ledger_entries" }

package model

import (
	"time"

	"github.com/google/uuid"
)

type ExchangeRate struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey"`
	BaseCurrency  string    `gorm:"type:varchar(10);not null"`
	QuoteCurrency string    `gorm:"type:varchar(10);not null"`
	Rate          float64   `gorm:"not null"`
	InverseRate   float64   `gorm:"not null"`
	Source        string    `gorm:"type:varchar(50)"`
	EffectiveDate time.Time `gorm:"type:date;not null"`
	CreatedAt     time.Time `gorm:"type:timestamptz;not null"`
}

func (ExchangeRate) TableName() string { return "exchange_rates" }

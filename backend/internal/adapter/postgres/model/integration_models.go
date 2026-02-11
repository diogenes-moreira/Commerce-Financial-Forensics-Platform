package model

import (
	"time"

	"github.com/google/uuid"
)

type Integration struct {
	ID           uuid.UUID  `gorm:"type:uuid;primaryKey"`
	Platform     string     `gorm:"type:varchar(50);not null"`
	Name         string     `gorm:"type:varchar(255);not null"`
	Config       string     `gorm:"type:text"` // JSON marshaled map[string]string
	Status       string     `gorm:"type:varchar(50);not null;default:'active'"`
	SyncStatus   string     `gorm:"type:varchar(50);not null;default:'idle'"`
	LastSyncedAt *time.Time `gorm:"type:timestamptz"`
	CreatedAt    time.Time  `gorm:"type:timestamptz;not null"`
	UpdatedAt    time.Time  `gorm:"type:timestamptz;not null"`
}

func (Integration) TableName() string { return "integrations" }

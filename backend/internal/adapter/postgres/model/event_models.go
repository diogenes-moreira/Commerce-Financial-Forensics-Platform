package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type ForensicEvent struct {
	ID            uuid.UUID      `gorm:"type:uuid;primaryKey"`
	EntityType    string         `gorm:"type:varchar(100);not null"`
	EntityID      uuid.UUID      `gorm:"type:uuid;not null"`
	EventType     string         `gorm:"type:varchar(100);not null"`
	EventTime     time.Time      `gorm:"type:timestamptz;not null"`
	EffectiveTime time.Time      `gorm:"type:timestamptz;not null"`
	Payload       datatypes.JSON `gorm:"type:jsonb;default:'{}'"`
	ActorID       string         `gorm:"type:varchar(255)"`
	HashIntegrity string         `gorm:"type:varchar(255);not null"`
	CreatedAt     time.Time      `gorm:"type:timestamptz;not null"`
}

func (ForensicEvent) TableName() string { return "forensic_events" }

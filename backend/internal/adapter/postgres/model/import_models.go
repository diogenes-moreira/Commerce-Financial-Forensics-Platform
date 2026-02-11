package model

import (
	"time"

	"github.com/google/uuid"
)

type ImportJob struct {
	ID            uuid.UUID  `gorm:"type:uuid;primaryKey"`
	Name          string     `gorm:"type:varchar(255);not null"`
	Source        string     `gorm:"type:varchar(50);not null"`
	EntityType    string     `gorm:"type:varchar(50);not null"`
	Status        string     `gorm:"type:varchar(50);not null;default:'pending'"`
	TotalRows     int        `gorm:"not null;default:0"`
	ProcessedRows int        `gorm:"not null;default:0"`
	FailedRows    int        `gorm:"not null;default:0"`
	ErrorLog      string     `gorm:"type:text"` // JSON array of strings
	SourceURI     string     `gorm:"type:text"`
	StartedAt     *time.Time `gorm:"type:timestamptz"`
	CompletedAt   *time.Time `gorm:"type:timestamptz"`
	CreatedAt     time.Time  `gorm:"type:timestamptz;not null"`
}

func (ImportJob) TableName() string { return "import_jobs" }

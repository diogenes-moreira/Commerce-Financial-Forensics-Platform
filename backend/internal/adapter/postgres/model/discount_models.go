package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type PromotionRule struct {
	ID            uuid.UUID      `gorm:"type:uuid;primaryKey"`
	Name          string         `gorm:"type:varchar(255);not null"`
	RuleType      string         `gorm:"type:varchar(50);not null"`
	Value         float64        `gorm:"not null"`
	Conditions    datatypes.JSON `gorm:"type:jsonb;default:'{}'"`
	FundingSource string         `gorm:"type:varchar(50);not null"`
	FundingPct    float64        `gorm:"not null;default:0"`
	MaxUsageCount int            `gorm:"not null;default:0"`
	CurrentUsage  int            `gorm:"not null;default:0"`
	ValidFrom     time.Time      `gorm:"type:timestamptz;not null"`
	ValidTo       time.Time      `gorm:"type:timestamptz;not null"`
	Status        string         `gorm:"type:varchar(50);not null;default:'active'"`
	Version       int            `gorm:"not null;default:1"`
	CreatedAt     time.Time      `gorm:"type:timestamptz;not null"`
	UpdatedAt     time.Time      `gorm:"type:timestamptz;not null"`
}

func (PromotionRule) TableName() string { return "promotion_rules" }

type DiscountApplication struct {
	ID                uuid.UUID `gorm:"type:uuid;primaryKey"`
	OrderID           uuid.UUID `gorm:"type:uuid;not null"`
	PromotionRuleID   uuid.UUID `gorm:"type:uuid;not null"`
	RuleVersion       int       `gorm:"not null"`
	DiscountCents     int64     `gorm:"not null"`
	Currency          string    `gorm:"type:varchar(10);not null;default:'USD'"`
	FundingSource     string    `gorm:"type:varchar(50);not null"`
	SellerShareCents  int64     `gorm:"not null;default:0"`
	PlatformShareCents int64   `gorm:"not null;default:0"`
	AppliedAt         time.Time `gorm:"type:timestamptz;not null"`
}

func (DiscountApplication) TableName() string { return "discount_applications" }

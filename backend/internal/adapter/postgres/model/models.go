package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type CloudAccount struct {
	ID           uuid.UUID  `gorm:"type:uuid;primaryKey"`
	Provider     string     `gorm:"type:varchar(50);not null"`
	Name         string     `gorm:"type:varchar(255);not null"`
	ExternalID   string     `gorm:"type:varchar(255);not null"`
	Status       string     `gorm:"type:varchar(50);not null;default:'active'"`
	SyncStatus   string     `gorm:"type:varchar(50);not null;default:'idle'"`
	LastSyncedAt *time.Time `gorm:"type:timestamptz"`
	CreatedAt    time.Time  `gorm:"type:timestamptz;not null"`
	UpdatedAt    time.Time  `gorm:"type:timestamptz;not null"`
}

func (CloudAccount) TableName() string { return "cloud_accounts" }

type CostRecord struct {
	ID             uuid.UUID      `gorm:"type:uuid;primaryKey"`
	CloudAccountID uuid.UUID      `gorm:"type:uuid;not null"`
	Service        string         `gorm:"type:varchar(255);not null"`
	Category       string         `gorm:"type:varchar(255);not null;default:'other'"`
	AmountCents    int64          `gorm:"not null"`
	Currency       string         `gorm:"type:varchar(10);not null;default:'USD'"`
	UsageDate      time.Time      `gorm:"type:date;not null"`
	Tags           datatypes.JSON `gorm:"type:jsonb;default:'{}'"`
	CreatedAt      time.Time      `gorm:"type:timestamptz;not null"`
}

func (CostRecord) TableName() string { return "cost_records" }

type Budget struct {
	ID               uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name             string    `gorm:"type:varchar(255);not null"`
	AmountCents      int64     `gorm:"not null"`
	Currency         string    `gorm:"type:varchar(10);not null;default:'USD'"`
	SpentCents       int64     `gorm:"not null;default:0"`
	PeriodStart      time.Time `gorm:"type:date;not null"`
	PeriodEnd        time.Time `gorm:"type:date;not null"`
	AlertThresholdPct int      `gorm:"not null;default:80"`
	Status           string    `gorm:"type:varchar(50);not null;default:'active'"`
	CreatedAt        time.Time `gorm:"type:timestamptz;not null"`
	UpdatedAt        time.Time `gorm:"type:timestamptz;not null"`
}

func (Budget) TableName() string { return "budgets" }

type BudgetAlert struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey"`
	BudgetID     uuid.UUID `gorm:"type:uuid;not null"`
	ThresholdPct int       `gorm:"not null"`
	ActualPct    int       `gorm:"not null"`
	Message      string    `gorm:"type:text;not null"`
	TriggeredAt  time.Time `gorm:"type:timestamptz;not null"`
}

func (BudgetAlert) TableName() string { return "budget_alerts" }

type Anomaly struct {
	ID             uuid.UUID  `gorm:"type:uuid;primaryKey"`
	CloudAccountID uuid.UUID  `gorm:"type:uuid;not null"`
	Service        string     `gorm:"type:varchar(255);not null"`
	ExpectedCents  int64      `gorm:"not null"`
	ActualCents    int64      `gorm:"not null"`
	Currency       string     `gorm:"type:varchar(10);not null;default:'USD'"`
	DeviationPct   float64    `gorm:"not null"`
	Severity       string     `gorm:"type:varchar(50);not null"`
	Status         string     `gorm:"type:varchar(50);not null;default:'open'"`
	DetectedAt     time.Time  `gorm:"type:timestamptz;not null"`
	ResolvedAt     *time.Time `gorm:"type:timestamptz"`
}

func (Anomaly) TableName() string { return "anomalies" }

type CostReport struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey"`
	Name        string         `gorm:"type:varchar(255);not null"`
	ReportType  string         `gorm:"type:varchar(50);not null"`
	PeriodStart time.Time      `gorm:"type:date;not null"`
	PeriodEnd   time.Time      `gorm:"type:date;not null"`
	TotalCents  int64          `gorm:"not null;default:0"`
	Currency    string         `gorm:"type:varchar(10);not null;default:'USD'"`
	Breakdown   datatypes.JSON `gorm:"type:jsonb;default:'{}'"`
	GeneratedAt time.Time      `gorm:"type:timestamptz;not null"`
}

func (CostReport) TableName() string { return "cost_reports" }

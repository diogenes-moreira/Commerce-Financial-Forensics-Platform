package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type Product struct {
	ID             uuid.UUID      `gorm:"type:uuid;primaryKey"`
	SKU            string         `gorm:"type:varchar(255);not null;uniqueIndex"`
	EAN            string         `gorm:"type:varchar(13)"`
	UPC            string         `gorm:"type:varchar(12)"`
	Name           string         `gorm:"type:varchar(255);not null"`
	Description    string         `gorm:"type:text"`
	Category       string         `gorm:"type:varchar(255)"`
	UnitCostCents  int64          `gorm:"not null"`
	UnitPriceCents int64          `gorm:"not null"`
	Currency       string         `gorm:"type:varchar(10);not null;default:'USD'"`
	Status         string         `gorm:"type:varchar(50);not null;default:'active'"`
	Metadata       datatypes.JSON `gorm:"type:jsonb;default:'{}'"`
	CreatedAt      time.Time      `gorm:"type:timestamptz;not null"`
	UpdatedAt      time.Time      `gorm:"type:timestamptz;not null"`
}

func (Product) TableName() string { return "products" }

type Seller struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey"`
	ExternalID    string    `gorm:"type:varchar(255)"`
	Code          string    `gorm:"type:varchar(255)"`
	Name          string    `gorm:"type:varchar(255);not null"`
	Email         string    `gorm:"type:varchar(255)"`
	CommissionPct float64   `gorm:"not null;default:0"`
	Status        string    `gorm:"type:varchar(50);not null;default:'active'"`
	CreatedAt     time.Time `gorm:"type:timestamptz;not null"`
	UpdatedAt     time.Time `gorm:"type:timestamptz;not null"`
}

func (Seller) TableName() string { return "sellers" }

type Customer struct {
	ID         uuid.UUID      `gorm:"type:uuid;primaryKey"`
	ExternalID string         `gorm:"type:varchar(255)"`
	Name       string         `gorm:"type:varchar(255);not null"`
	Email      string         `gorm:"type:varchar(255)"`
	Segment    string         `gorm:"type:varchar(100)"`
	Metadata   datatypes.JSON `gorm:"type:jsonb;default:'{}'"`
	CreatedAt  time.Time      `gorm:"type:timestamptz;not null"`
	UpdatedAt  time.Time      `gorm:"type:timestamptz;not null"`
}

func (Customer) TableName() string { return "customers" }

type Order struct {
	ID            uuid.UUID      `gorm:"type:uuid;primaryKey"`
	ExternalID    string         `gorm:"type:varchar(255);not null"`
	SellerID      uuid.UUID      `gorm:"type:uuid;not null"`
	CustomerID    uuid.UUID      `gorm:"type:uuid;not null"`
	Status        string         `gorm:"type:varchar(50);not null;default:'pending'"`
	SubtotalCents int64          `gorm:"not null;default:0"`
	DiscountCents int64          `gorm:"not null;default:0"`
	ShippingCents int64          `gorm:"not null;default:0"`
	TaxCents      int64          `gorm:"not null;default:0"`
	TotalCents    int64          `gorm:"not null;default:0"`
	Currency      string         `gorm:"type:varchar(10);not null;default:'USD'"`
	OrderDate     time.Time      `gorm:"type:timestamptz;not null"`
	Metadata      datatypes.JSON `gorm:"type:jsonb;default:'{}'"`
	CreatedAt     time.Time      `gorm:"type:timestamptz;not null"`
	UpdatedAt     time.Time      `gorm:"type:timestamptz;not null"`
}

func (Order) TableName() string { return "orders" }

type OrderItem struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey"`
	OrderID        uuid.UUID `gorm:"type:uuid;not null"`
	ProductID      uuid.UUID `gorm:"type:uuid;not null"`
	SKU            string    `gorm:"type:varchar(255);not null"`
	ProductName    string    `gorm:"type:varchar(255);not null"`
	Quantity       int       `gorm:"not null"`
	UnitPriceCents int64     `gorm:"not null"`
	UnitCostCents  int64     `gorm:"not null"`
	DiscountCents  int64     `gorm:"not null;default:0"`
	TotalCents     int64     `gorm:"not null"`
	Currency       string    `gorm:"type:varchar(10);not null;default:'USD'"`
	CreatedAt      time.Time `gorm:"type:timestamptz;not null"`
}

func (OrderItem) TableName() string { return "order_items" }

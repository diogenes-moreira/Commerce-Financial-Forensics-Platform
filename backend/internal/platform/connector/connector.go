package connector

import (
	"context"
	"time"
)

type RawProduct struct {
	ExternalID  string
	SKU         string
	Name        string
	Description string
	Category    string
	CostCents   int64
	PriceCents  int64
	Currency    string
}

type RawOrder struct {
	ExternalID    string
	SellerID      string
	CustomerID    string
	Items         []RawOrderItem
	SubtotalCents int64
	DiscountCents int64
	ShippingCents int64
	TaxCents      int64
	TotalCents    int64
	Currency      string
	OrderDate     time.Time
}

type RawOrderItem struct {
	SKU            string
	ProductName    string
	Quantity       int
	UnitPriceCents int64
	UnitCostCents  int64
}

type RawCustomer struct {
	ExternalID string
	Name       string
	Email      string
	Segment    string
}

type Connector interface {
	FetchProducts(ctx context.Context, since time.Time) ([]RawProduct, error)
	FetchOrders(ctx context.Context, since time.Time) ([]RawOrder, error)
	FetchCustomers(ctx context.Context, since time.Time) ([]RawCustomer, error)
	HandleWebhook(ctx context.Context, eventType string, payload []byte) error
}

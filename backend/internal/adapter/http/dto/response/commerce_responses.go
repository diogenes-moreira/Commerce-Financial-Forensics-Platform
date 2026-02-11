package response

import (
	"time"

	"github.com/diogenes/costforensics/backend/internal/domain/customer"
	"github.com/diogenes/costforensics/backend/internal/domain/order"
	"github.com/diogenes/costforensics/backend/internal/domain/product"
	"github.com/diogenes/costforensics/backend/internal/domain/seller"
)

type Product struct {
	ID             string            `json:"id"`
	SKU            string            `json:"sku"`
	EAN            string            `json:"ean"`
	UPC            string            `json:"upc"`
	Name           string            `json:"name"`
	Description    string            `json:"description"`
	Category       string            `json:"category"`
	UnitCostCents  int64             `json:"unit_cost_cents"`
	UnitPriceCents int64             `json:"unit_price_cents"`
	Currency       string            `json:"currency"`
	GrossMarginPct float64           `json:"gross_margin_pct"`
	Status         string            `json:"status"`
	Metadata       map[string]string `json:"metadata"`
	CreatedAt      time.Time         `json:"created_at"`
	UpdatedAt      time.Time         `json:"updated_at"`
}

func ProductFromDomain(p *product.Product) Product {
	return Product{
		ID: p.ID().String(), SKU: p.SKU(), EAN: p.EAN(), UPC: p.UPC(), Name: p.Name(),
		Description: p.Description(), Category: p.Category(),
		UnitCostCents: p.UnitCostCents().Amount(),
		UnitPriceCents: p.UnitPriceCents().Amount(),
		Currency:       p.UnitPriceCents().Currency(),
		GrossMarginPct: p.GrossMarginPct(),
		Status:   string(p.Status()), Metadata: p.Metadata(),
		CreatedAt: p.CreatedAt(), UpdatedAt: p.UpdatedAt(),
	}
}

type Seller struct {
	ID            string    `json:"id"`
	ExternalID    string    `json:"external_id"`
	Code          string    `json:"code"`
	Name          string    `json:"name"`
	Email         string    `json:"email"`
	CommissionPct float64   `json:"commission_pct"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func SellerFromDomain(s *seller.Seller) Seller {
	return Seller{
		ID: s.ID().String(), ExternalID: s.ExternalID(),
		Code: s.Code(), Name: s.Name(), Email: s.Email(),
		CommissionPct: s.CommissionPct(), Status: string(s.Status()),
		CreatedAt: s.CreatedAt(), UpdatedAt: s.UpdatedAt(),
	}
}

type Customer struct {
	ID         string            `json:"id"`
	ExternalID string            `json:"external_id"`
	Name       string            `json:"name"`
	Email      string            `json:"email"`
	Segment    string            `json:"segment"`
	Metadata   map[string]string `json:"metadata"`
	CreatedAt  time.Time         `json:"created_at"`
	UpdatedAt  time.Time         `json:"updated_at"`
}

func CustomerFromDomain(c *customer.Customer) Customer {
	return Customer{
		ID: c.ID().String(), ExternalID: c.ExternalID(),
		Name: c.Name(), Email: c.Email(), Segment: c.Segment(),
		Metadata: c.Metadata(), CreatedAt: c.CreatedAt(), UpdatedAt: c.UpdatedAt(),
	}
}

type OrderItem struct {
	ID             string    `json:"id"`
	OrderID        string    `json:"order_id"`
	ProductID      string    `json:"product_id"`
	SKU            string    `json:"sku"`
	ProductName    string    `json:"product_name"`
	Quantity       int       `json:"quantity"`
	UnitPriceCents int64     `json:"unit_price_cents"`
	UnitCostCents  int64     `json:"unit_cost_cents"`
	DiscountCents  int64     `json:"discount_cents"`
	TotalCents     int64     `json:"total_cents"`
	CreatedAt      time.Time `json:"created_at"`
}

func OrderItemFromDomain(i *order.OrderItem) OrderItem {
	return OrderItem{
		ID: i.ID().String(), OrderID: i.OrderID().String(),
		ProductID: i.ProductID().String(), SKU: i.SKU(),
		ProductName: i.ProductName(), Quantity: i.Quantity(),
		UnitPriceCents: i.UnitPriceCents().Amount(),
		UnitCostCents:  i.UnitCostCents().Amount(),
		DiscountCents:  i.DiscountCents().Amount(),
		TotalCents:     i.TotalCents().Amount(),
		CreatedAt:      i.CreatedAt(),
	}
}

type Order struct {
	ID            string            `json:"id"`
	ExternalID    string            `json:"external_id"`
	SellerID      string            `json:"seller_id"`
	CustomerID    string            `json:"customer_id"`
	Status        string            `json:"status"`
	SubtotalCents int64             `json:"subtotal_cents"`
	DiscountCents int64             `json:"discount_cents"`
	ShippingCents int64             `json:"shipping_cents"`
	TaxCents      int64             `json:"tax_cents"`
	TotalCents    int64             `json:"total_cents"`
	Currency      string            `json:"currency"`
	OrderDate     string            `json:"order_date"`
	Items         []OrderItem       `json:"items"`
	Metadata      map[string]string `json:"metadata"`
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at"`
}

func OrderFromDomain(o *order.Order) Order {
	items := make([]OrderItem, len(o.Items()))
	for i, item := range o.Items() {
		items[i] = OrderItemFromDomain(item)
	}
	return Order{
		ID: o.ID().String(), ExternalID: o.ExternalID(),
		SellerID: o.SellerID().String(), CustomerID: o.CustomerID().String(),
		Status:        string(o.Status()),
		SubtotalCents: o.SubtotalCents().Amount(),
		DiscountCents: o.DiscountCents().Amount(),
		ShippingCents: o.ShippingCents().Amount(),
		TaxCents:      o.TaxCents().Amount(),
		TotalCents:    o.TotalCents().Amount(),
		Currency:      o.Currency(),
		OrderDate:     o.OrderDate().Format("2006-01-02"),
		Items:         items, Metadata: o.Metadata(),
		CreatedAt: o.CreatedAt(), UpdatedAt: o.UpdatedAt(),
	}
}

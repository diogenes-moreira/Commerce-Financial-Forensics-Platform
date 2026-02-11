package mapper

import (
	"encoding/json"

	"github.com/diogenes/costforensics/backend/internal/adapter/postgres/model"
	"github.com/diogenes/costforensics/backend/internal/domain/customer"
	"github.com/diogenes/costforensics/backend/internal/domain/order"
	"github.com/diogenes/costforensics/backend/internal/domain/product"
	"github.com/diogenes/costforensics/backend/internal/domain/seller"
	"github.com/diogenes/costforensics/backend/internal/domain/shared"
)

// Product

func ProductToModel(p *product.Product) *model.Product {
	metaJSON, _ := json.Marshal(p.Metadata())
	return &model.Product{
		ID:             p.ID(),
		SKU:            p.SKU(),
		EAN:            p.EAN(),
		UPC:            p.UPC(),
		Name:           p.Name(),
		Description:    p.Description(),
		Category:       p.Category(),
		UnitCostCents:  p.UnitCostCents().Amount(),
		UnitPriceCents: p.UnitPriceCents().Amount(),
		Currency:       p.UnitPriceCents().Currency(),
		Status:         string(p.Status()),
		Metadata:       metaJSON,
		CreatedAt:      p.CreatedAt(),
		UpdatedAt:      p.UpdatedAt(),
	}
}

func ProductToDomain(m *model.Product) *product.Product {
	unitCost := shared.MustNewMoney(m.UnitCostCents, m.Currency)
	unitPrice := shared.MustNewMoney(m.UnitPriceCents, m.Currency)
	var meta map[string]string
	_ = json.Unmarshal(m.Metadata, &meta)
	if meta == nil {
		meta = make(map[string]string)
	}
	return product.HydrateProduct(
		m.ID, m.SKU, m.EAN, m.UPC, m.Name, m.Description, m.Category,
		unitCost, unitPrice, product.Status(m.Status),
		meta, m.CreatedAt, m.UpdatedAt,
	)
}

// Seller

func SellerToModel(s *seller.Seller) *model.Seller {
	return &model.Seller{
		ID:            s.ID(),
		ExternalID:    s.ExternalID(),
		Code:          s.Code(),
		Name:          s.Name(),
		Email:         s.Email(),
		CommissionPct: s.CommissionPct(),
		Status:        string(s.Status()),
		CreatedAt:     s.CreatedAt(),
		UpdatedAt:     s.UpdatedAt(),
	}
}

func SellerToDomain(m *model.Seller) *seller.Seller {
	return seller.HydrateSeller(
		m.ID, m.ExternalID, m.Code, m.Name, m.Email,
		m.CommissionPct, seller.Status(m.Status),
		m.CreatedAt, m.UpdatedAt,
	)
}

// Customer

func CustomerToModel(c *customer.Customer) *model.Customer {
	metaJSON, _ := json.Marshal(c.Metadata())
	return &model.Customer{
		ID:         c.ID(),
		ExternalID: c.ExternalID(),
		Name:       c.Name(),
		Email:      c.Email(),
		Segment:    c.Segment(),
		Metadata:   metaJSON,
		CreatedAt:  c.CreatedAt(),
		UpdatedAt:  c.UpdatedAt(),
	}
}

func CustomerToDomain(m *model.Customer) *customer.Customer {
	var meta map[string]string
	_ = json.Unmarshal(m.Metadata, &meta)
	if meta == nil {
		meta = make(map[string]string)
	}
	return customer.HydrateCustomer(
		m.ID, m.ExternalID, m.Name, m.Email,
		m.Segment, meta, m.CreatedAt, m.UpdatedAt,
	)
}

// Order

func OrderToModel(o *order.Order) *model.Order {
	metaJSON, _ := json.Marshal(o.Metadata())
	return &model.Order{
		ID:            o.ID(),
		ExternalID:    o.ExternalID(),
		SellerID:      o.SellerID(),
		CustomerID:    o.CustomerID(),
		Status:        string(o.Status()),
		SubtotalCents: o.SubtotalCents().Amount(),
		DiscountCents: o.DiscountCents().Amount(),
		ShippingCents: o.ShippingCents().Amount(),
		TaxCents:      o.TaxCents().Amount(),
		TotalCents:    o.TotalCents().Amount(),
		Currency:      o.Currency(),
		OrderDate:     o.OrderDate(),
		Metadata:      metaJSON,
		CreatedAt:     o.CreatedAt(),
		UpdatedAt:     o.UpdatedAt(),
	}
}

func OrderToDomain(m *model.Order, items []*order.OrderItem) *order.Order {
	var meta map[string]string
	_ = json.Unmarshal(m.Metadata, &meta)
	if meta == nil {
		meta = make(map[string]string)
	}
	subtotal := shared.MustNewMoney(m.SubtotalCents, m.Currency)
	discount := shared.MustNewMoney(m.DiscountCents, m.Currency)
	shipping := shared.MustNewMoney(m.ShippingCents, m.Currency)
	tax := shared.MustNewMoney(m.TaxCents, m.Currency)
	total := shared.MustNewMoney(m.TotalCents, m.Currency)
	return order.HydrateOrder(
		m.ID, m.ExternalID, m.SellerID, m.CustomerID,
		order.Status(m.Status), subtotal, discount, shipping, tax, total,
		m.Currency, m.OrderDate, items, meta, m.CreatedAt, m.UpdatedAt,
	)
}

// OrderItem

func OrderItemToModel(i *order.OrderItem) *model.OrderItem {
	return &model.OrderItem{
		ID:             i.ID(),
		OrderID:        i.OrderID(),
		ProductID:      i.ProductID(),
		SKU:            i.SKU(),
		ProductName:    i.ProductName(),
		Quantity:       i.Quantity(),
		UnitPriceCents: i.UnitPriceCents().Amount(),
		UnitCostCents:  i.UnitCostCents().Amount(),
		DiscountCents:  i.DiscountCents().Amount(),
		TotalCents:     i.TotalCents().Amount(),
		Currency:       i.UnitPriceCents().Currency(),
		CreatedAt:      i.CreatedAt(),
	}
}

func OrderItemToDomain(m *model.OrderItem) *order.OrderItem {
	unitPrice := shared.MustNewMoney(m.UnitPriceCents, m.Currency)
	unitCost := shared.MustNewMoney(m.UnitCostCents, m.Currency)
	discount := shared.MustNewMoney(m.DiscountCents, m.Currency)
	total := shared.MustNewMoney(m.TotalCents, m.Currency)
	return order.HydrateOrderItem(
		m.ID, m.OrderID, m.ProductID,
		m.SKU, m.ProductName, m.Quantity,
		unitPrice, unitCost, discount, total,
		m.CreatedAt,
	)
}

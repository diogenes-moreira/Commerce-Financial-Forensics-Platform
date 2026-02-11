package product

import (
	"time"

	"github.com/diogenes/costforensics/backend/internal/domain/shared"
	"github.com/google/uuid"
)

type Status string

const (
	StatusActive       Status = "active"
	StatusInactive     Status = "inactive"
	StatusDiscontinued Status = "discontinued"
)

type Product struct {
	id             uuid.UUID
	sku            string
	ean            string
	upc            string
	name           string
	description    string
	category       string
	unitCostCents  shared.Money
	unitPriceCents shared.Money
	status         Status
	metadata       map[string]string
	createdAt      time.Time
	updatedAt      time.Time
}

func NewProduct(sku, name string, unitCost, unitPrice shared.Money) (*Product, error) {
	if sku == "" {
		return nil, ErrSKUEmpty
	}
	if name == "" {
		return nil, ErrNameEmpty
	}
	if unitCost.GreaterThan(unitPrice) {
		return nil, ErrPriceBelowCost
	}

	now := time.Now().UTC()
	return &Product{
		id:             uuid.New(),
		sku:            sku,
		name:           name,
		unitCostCents:  unitCost,
		unitPriceCents: unitPrice,
		status:         StatusActive,
		metadata:       make(map[string]string),
		createdAt:      now,
		updatedAt:      now,
	}, nil
}

func HydrateProduct(
	id uuid.UUID, sku, ean, upc, name, description, category string,
	unitCost, unitPrice shared.Money, status Status,
	metadata map[string]string, createdAt, updatedAt time.Time,
) *Product {
	if metadata == nil {
		metadata = make(map[string]string)
	}
	return &Product{
		id: id, sku: sku, ean: ean, upc: upc,
		name: name, description: description,
		category: category, unitCostCents: unitCost, unitPriceCents: unitPrice,
		status: status, metadata: metadata,
		createdAt: createdAt, updatedAt: updatedAt,
	}
}

func (p *Product) ID() uuid.UUID            { return p.id }
func (p *Product) SKU() string              { return p.sku }
func (p *Product) EAN() string              { return p.ean }
func (p *Product) UPC() string              { return p.upc }
func (p *Product) Name() string             { return p.name }
func (p *Product) Description() string      { return p.description }
func (p *Product) Category() string         { return p.category }
func (p *Product) UnitCostCents() shared.Money  { return p.unitCostCents }
func (p *Product) UnitPriceCents() shared.Money { return p.unitPriceCents }
func (p *Product) Status() Status           { return p.status }
func (p *Product) Metadata() map[string]string  { return p.metadata }
func (p *Product) CreatedAt() time.Time     { return p.createdAt }
func (p *Product) UpdatedAt() time.Time     { return p.updatedAt }

func (p *Product) UpdatePrice(newPrice shared.Money) error {
	if p.unitCostCents.GreaterThan(newPrice) {
		return ErrPriceBelowCost
	}
	p.unitPriceCents = newPrice
	p.updatedAt = time.Now().UTC()
	return nil
}

func (p *Product) UpdateDetails(name, description, category string) error {
	if name == "" {
		return ErrNameEmpty
	}
	p.name = name
	p.description = description
	p.category = category
	p.updatedAt = time.Now().UTC()
	return nil
}

func (p *Product) Deactivate() {
	p.status = StatusInactive
	p.updatedAt = time.Now().UTC()
}

func (p *Product) GrossMarginPct() float64 {
	if p.unitPriceCents.IsZero() {
		return 0
	}
	margin := float64(p.unitPriceCents.Amount()-p.unitCostCents.Amount()) / float64(p.unitPriceCents.Amount()) * 100
	return margin
}

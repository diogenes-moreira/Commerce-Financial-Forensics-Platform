package order

import (
	"time"

	"github.com/diogenes/costforensics/backend/internal/domain/shared"
	"github.com/google/uuid"
)

type OrderItem struct {
	id             uuid.UUID
	orderID        uuid.UUID
	productID      uuid.UUID
	sku            string
	productName    string
	quantity       int
	unitPriceCents shared.Money
	unitCostCents  shared.Money
	discountCents  shared.Money
	totalCents     shared.Money
	createdAt      time.Time
}

func newOrderItem(
	orderID, productID uuid.UUID,
	sku, productName string,
	qty int, unitPrice, unitCost shared.Money,
) (*OrderItem, error) {
	if qty <= 0 {
		return nil, ErrInvalidQuantity
	}
	lineTotal := unitPrice.Multiply(float64(qty))
	return &OrderItem{
		id:             uuid.New(),
		orderID:        orderID,
		productID:      productID,
		sku:            sku,
		productName:    productName,
		quantity:       qty,
		unitPriceCents: unitPrice,
		unitCostCents:  unitCost,
		discountCents:  shared.ZeroMoney(unitPrice.Currency()),
		totalCents:     lineTotal,
		createdAt:      time.Now().UTC(),
	}, nil
}

func HydrateOrderItem(
	id, orderID, productID uuid.UUID,
	sku, productName string, quantity int,
	unitPrice, unitCost, discount, total shared.Money,
	createdAt time.Time,
) *OrderItem {
	return &OrderItem{
		id: id, orderID: orderID, productID: productID,
		sku: sku, productName: productName, quantity: quantity,
		unitPriceCents: unitPrice, unitCostCents: unitCost,
		discountCents: discount, totalCents: total,
		createdAt: createdAt,
	}
}

func (i *OrderItem) ID() uuid.UUID            { return i.id }
func (i *OrderItem) OrderID() uuid.UUID        { return i.orderID }
func (i *OrderItem) ProductID() uuid.UUID      { return i.productID }
func (i *OrderItem) SKU() string               { return i.sku }
func (i *OrderItem) ProductName() string       { return i.productName }
func (i *OrderItem) Quantity() int             { return i.quantity }
func (i *OrderItem) UnitPriceCents() shared.Money { return i.unitPriceCents }
func (i *OrderItem) UnitCostCents() shared.Money  { return i.unitCostCents }
func (i *OrderItem) DiscountCents() shared.Money   { return i.discountCents }
func (i *OrderItem) TotalCents() shared.Money      { return i.totalCents }
func (i *OrderItem) CreatedAt() time.Time      { return i.createdAt }

func (i *OrderItem) COGSCents() shared.Money {
	return i.unitCostCents.Multiply(float64(i.quantity))
}

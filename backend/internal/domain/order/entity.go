package order

import (
	"time"

	"github.com/diogenes/costforensics/backend/internal/domain/shared"
	"github.com/google/uuid"
)

type Status string

const (
	StatusPending   Status = "pending"
	StatusConfirmed Status = "confirmed"
	StatusShipped   Status = "shipped"
	StatusDelivered Status = "delivered"
	StatusCancelled Status = "cancelled"
	StatusRefunded  Status = "refunded"
)

type Order struct {
	id            uuid.UUID
	externalID    string
	sellerID      uuid.UUID
	customerID    uuid.UUID
	status        Status
	subtotalCents shared.Money
	discountCents shared.Money
	shippingCents shared.Money
	taxCents      shared.Money
	totalCents    shared.Money
	currency      string
	orderDate     time.Time
	items         []*OrderItem
	metadata      map[string]string
	createdAt     time.Time
	updatedAt     time.Time
}

func NewOrder(externalID string, sellerID, customerID uuid.UUID, currency string) (*Order, error) {
	if externalID == "" {
		return nil, ErrExternalIDEmpty
	}
	if currency == "" {
		currency = "USD"
	}

	now := time.Now().UTC()
	return &Order{
		id:            uuid.New(),
		externalID:    externalID,
		sellerID:      sellerID,
		customerID:    customerID,
		status:        StatusPending,
		subtotalCents: shared.ZeroMoney(currency),
		discountCents: shared.ZeroMoney(currency),
		shippingCents: shared.ZeroMoney(currency),
		taxCents:      shared.ZeroMoney(currency),
		totalCents:    shared.ZeroMoney(currency),
		currency:      currency,
		orderDate:     now,
		items:         make([]*OrderItem, 0),
		metadata:      make(map[string]string),
		createdAt:     now,
		updatedAt:     now,
	}, nil
}

func HydrateOrder(
	id uuid.UUID, externalID string, sellerID, customerID uuid.UUID,
	status Status, subtotal, discount, shipping, tax, total shared.Money,
	currency string, orderDate time.Time, items []*OrderItem,
	metadata map[string]string, createdAt, updatedAt time.Time,
) *Order {
	if items == nil {
		items = make([]*OrderItem, 0)
	}
	if metadata == nil {
		metadata = make(map[string]string)
	}
	return &Order{
		id: id, externalID: externalID, sellerID: sellerID,
		customerID: customerID, status: status,
		subtotalCents: subtotal, discountCents: discount,
		shippingCents: shipping, taxCents: tax, totalCents: total,
		currency: currency, orderDate: orderDate, items: items,
		metadata: metadata, createdAt: createdAt, updatedAt: updatedAt,
	}
}

func (o *Order) ID() uuid.UUID             { return o.id }
func (o *Order) ExternalID() string         { return o.externalID }
func (o *Order) SellerID() uuid.UUID        { return o.sellerID }
func (o *Order) CustomerID() uuid.UUID      { return o.customerID }
func (o *Order) Status() Status             { return o.status }
func (o *Order) SubtotalCents() shared.Money { return o.subtotalCents }
func (o *Order) DiscountCents() shared.Money { return o.discountCents }
func (o *Order) ShippingCents() shared.Money { return o.shippingCents }
func (o *Order) TaxCents() shared.Money      { return o.taxCents }
func (o *Order) TotalCents() shared.Money    { return o.totalCents }
func (o *Order) Currency() string            { return o.currency }
func (o *Order) OrderDate() time.Time        { return o.orderDate }
func (o *Order) Items() []*OrderItem         { return o.items }
func (o *Order) Metadata() map[string]string { return o.metadata }
func (o *Order) CreatedAt() time.Time        { return o.createdAt }
func (o *Order) UpdatedAt() time.Time        { return o.updatedAt }

func (o *Order) AddItem(productID uuid.UUID, sku, name string, qty int, unitPrice, unitCost shared.Money) error {
	item, err := newOrderItem(o.id, productID, sku, name, qty, unitPrice, unitCost)
	if err != nil {
		return err
	}
	o.items = append(o.items, item)
	o.recalculate()
	return nil
}

func (o *Order) ApplyDiscount(amount shared.Money) {
	o.discountCents = amount
	o.recalculate()
}

func (o *Order) SetShipping(amount shared.Money) {
	o.shippingCents = amount
	o.recalculate()
}

func (o *Order) SetTax(amount shared.Money) {
	o.taxCents = amount
	o.recalculate()
}

func (o *Order) recalculate() {
	subtotal := shared.ZeroMoney(o.currency)
	for _, item := range o.items {
		subtotal, _ = subtotal.Add(item.totalCents)
	}
	o.subtotalCents = subtotal

	// total = subtotal - discount + shipping + tax
	total := subtotal
	if o.discountCents.Amount() > 0 {
		result, err := total.Subtract(o.discountCents)
		if err == nil {
			total = result
		}
	}
	total, _ = total.Add(o.shippingCents)
	total, _ = total.Add(o.taxCents)
	o.totalCents = total
	o.updatedAt = time.Now().UTC()
}

func (o *Order) Confirm() error {
	if o.status != StatusPending {
		return ErrInvalidTransition
	}
	o.status = StatusConfirmed
	o.updatedAt = time.Now().UTC()
	return nil
}

func (o *Order) Ship() error {
	if o.status != StatusConfirmed {
		return ErrInvalidTransition
	}
	o.status = StatusShipped
	o.updatedAt = time.Now().UTC()
	return nil
}

func (o *Order) Deliver() error {
	if o.status != StatusShipped {
		return ErrInvalidTransition
	}
	o.status = StatusDelivered
	o.updatedAt = time.Now().UTC()
	return nil
}

func (o *Order) Cancel() error {
	if o.status == StatusDelivered || o.status == StatusRefunded {
		return ErrInvalidTransition
	}
	o.status = StatusCancelled
	o.updatedAt = time.Now().UTC()
	return nil
}

func (o *Order) Refund() error {
	if o.status != StatusDelivered && o.status != StatusConfirmed && o.status != StatusShipped {
		return ErrInvalidTransition
	}
	o.status = StatusRefunded
	o.updatedAt = time.Now().UTC()
	return nil
}

func (o *Order) GrossMarginCents() shared.Money {
	totalCOGS := shared.ZeroMoney(o.currency)
	for _, item := range o.items {
		totalCOGS, _ = totalCOGS.Add(item.COGSCents())
	}
	result, err := o.subtotalCents.Subtract(totalCOGS)
	if err != nil {
		return shared.ZeroMoney(o.currency)
	}
	return result
}

func (o *Order) NetMarginCents() shared.Money {
	totalCOGS := shared.ZeroMoney(o.currency)
	for _, item := range o.items {
		totalCOGS, _ = totalCOGS.Add(item.COGSCents())
	}
	// net = total - COGS
	result, err := o.totalCents.Subtract(totalCOGS)
	if err != nil {
		return shared.ZeroMoney(o.currency)
	}
	return result
}

package payment

import (
	"time"

	"github.com/diogenes/costforensics/backend/internal/domain/shared"
	"github.com/google/uuid"
)

// Direction indicates money flow: inbound = from customer, outbound = to seller
type Direction string

const (
	DirectionInbound  Direction = "inbound"  // customer payment
	DirectionOutbound Direction = "outbound" // seller payout/disbursement
)

type Status string

const (
	StatusPending   Status = "pending"
	StatusProcessed Status = "processed"
	StatusFailed    Status = "failed"
	StatusRefunded  Status = "refunded"
)

type Payment struct {
	id              uuid.UUID
	orderID         uuid.UUID
	direction       Direction
	counterpartyID  uuid.UUID // customer ID (inbound) or seller ID (outbound)
	amountCents     shared.Money
	method          string // "credit_card", "bank_transfer", "wallet", "commission_payout", etc.
	externalRef     string // payment gateway reference
	status          Status
	processedAt     *time.Time
	failureReason   string
	metadata        map[string]string
	createdAt       time.Time
	updatedAt       time.Time
}

func NewPayment(
	orderID uuid.UUID, direction Direction, counterpartyID uuid.UUID,
	amount shared.Money, method, externalRef string,
) (*Payment, error) {
	if direction != DirectionInbound && direction != DirectionOutbound {
		return nil, ErrInvalidDirection
	}
	if amount.IsZero() {
		return nil, ErrAmountRequired
	}

	now := time.Now().UTC()
	return &Payment{
		id:             uuid.New(),
		orderID:        orderID,
		direction:      direction,
		counterpartyID: counterpartyID,
		amountCents:    amount,
		method:         method,
		externalRef:    externalRef,
		status:         StatusPending,
		metadata:       make(map[string]string),
		createdAt:      now,
		updatedAt:      now,
	}, nil
}

func HydratePayment(
	id, orderID uuid.UUID, direction Direction, counterpartyID uuid.UUID,
	amount shared.Money, method, externalRef string, status Status,
	processedAt *time.Time, failureReason string,
	metadata map[string]string, createdAt, updatedAt time.Time,
) *Payment {
	if metadata == nil {
		metadata = make(map[string]string)
	}
	return &Payment{
		id: id, orderID: orderID, direction: direction,
		counterpartyID: counterpartyID, amountCents: amount,
		method: method, externalRef: externalRef, status: status,
		processedAt: processedAt, failureReason: failureReason,
		metadata: metadata, createdAt: createdAt, updatedAt: updatedAt,
	}
}

func (p *Payment) ID() uuid.UUID             { return p.id }
func (p *Payment) OrderID() uuid.UUID         { return p.orderID }
func (p *Payment) Direction() Direction       { return p.direction }
func (p *Payment) CounterpartyID() uuid.UUID  { return p.counterpartyID }
func (p *Payment) AmountCents() shared.Money  { return p.amountCents }
func (p *Payment) Method() string             { return p.method }
func (p *Payment) ExternalRef() string        { return p.externalRef }
func (p *Payment) Status() Status             { return p.status }
func (p *Payment) ProcessedAt() *time.Time    { return p.processedAt }
func (p *Payment) FailureReason() string      { return p.failureReason }
func (p *Payment) Metadata() map[string]string { return p.metadata }
func (p *Payment) CreatedAt() time.Time       { return p.createdAt }
func (p *Payment) UpdatedAt() time.Time       { return p.updatedAt }

func (p *Payment) MarkProcessed() error {
	if p.status != StatusPending {
		return ErrInvalidStatus
	}
	now := time.Now().UTC()
	p.status = StatusProcessed
	p.processedAt = &now
	p.updatedAt = now
	return nil
}

func (p *Payment) MarkFailed(reason string) error {
	if p.status != StatusPending {
		return ErrInvalidStatus
	}
	p.status = StatusFailed
	p.failureReason = reason
	p.updatedAt = time.Now().UTC()
	return nil
}

func (p *Payment) MarkRefunded() error {
	if p.status != StatusProcessed {
		return ErrInvalidStatus
	}
	p.status = StatusRefunded
	p.updatedAt = time.Now().UTC()
	return nil
}

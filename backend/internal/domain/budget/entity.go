package budget

import (
	"time"

	"github.com/diogenes/costforensics/backend/internal/domain/shared"
	"github.com/google/uuid"
)

type Status string

const (
	StatusActive  Status = "active"
	StatusClosed  Status = "closed"
)

type Budget struct {
	id               uuid.UUID
	name             string
	amount           shared.Money
	spent            shared.Money
	periodStart      time.Time
	periodEnd        time.Time
	alertThresholdPct int
	status           Status
	createdAt        time.Time
	updatedAt        time.Time
}

func NewBudget(name string, amount shared.Money, periodStart, periodEnd time.Time, alertThresholdPct int) (*Budget, error) {
	if name == "" {
		return nil, ErrBudgetNameEmpty
	}
	if !periodStart.Before(periodEnd) {
		return nil, ErrInvalidPeriod
	}
	if alertThresholdPct <= 0 || alertThresholdPct > 100 {
		alertThresholdPct = 80
	}

	now := time.Now().UTC()
	return &Budget{
		id:               uuid.New(),
		name:             name,
		amount:           amount,
		spent:            shared.ZeroMoney(amount.Currency()),
		periodStart:      periodStart,
		periodEnd:        periodEnd,
		alertThresholdPct: alertThresholdPct,
		status:           StatusActive,
		createdAt:        now,
		updatedAt:        now,
	}, nil
}

func HydrateBudget(
	id uuid.UUID, name string, amount, spent shared.Money,
	periodStart, periodEnd time.Time, alertThresholdPct int,
	status Status, createdAt, updatedAt time.Time,
) *Budget {
	return &Budget{
		id: id, name: name, amount: amount, spent: spent,
		periodStart: periodStart, periodEnd: periodEnd,
		alertThresholdPct: alertThresholdPct, status: status,
		createdAt: createdAt, updatedAt: updatedAt,
	}
}

func (b *Budget) ID() uuid.UUID            { return b.id }
func (b *Budget) Name() string              { return b.name }
func (b *Budget) Amount() shared.Money      { return b.amount }
func (b *Budget) Spent() shared.Money       { return b.spent }
func (b *Budget) PeriodStart() time.Time    { return b.periodStart }
func (b *Budget) PeriodEnd() time.Time      { return b.periodEnd }
func (b *Budget) AlertThresholdPct() int    { return b.alertThresholdPct }
func (b *Budget) Status() Status            { return b.status }
func (b *Budget) CreatedAt() time.Time      { return b.createdAt }
func (b *Budget) UpdatedAt() time.Time      { return b.updatedAt }

func (b *Budget) UsagePct() int {
	if b.amount.IsZero() {
		return 0
	}
	return int(float64(b.spent.Amount()) / float64(b.amount.Amount()) * 100)
}

// RecordSpend adds spending and returns an alert if threshold is crossed.
func (b *Budget) RecordSpend(amount shared.Money) (*Alert, error) {
	newSpent, err := b.spent.Add(amount)
	if err != nil {
		return nil, err
	}
	b.spent = newSpent
	b.updatedAt = time.Now().UTC()

	// Check if we crossed the alert threshold
	usagePct := b.UsagePct()
	if usagePct >= b.alertThresholdPct {
		return newAlert(b.id, b.alertThresholdPct, usagePct), nil
	}
	return nil, nil
}

func (b *Budget) Close() {
	b.status = StatusClosed
	b.updatedAt = time.Now().UTC()
}

func (b *Budget) Rename(name string) error {
	if name == "" {
		return ErrBudgetNameEmpty
	}
	b.name = name
	b.updatedAt = time.Now().UTC()
	return nil
}

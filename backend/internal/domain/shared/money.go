package shared

import (
	"errors"
	"fmt"
	"math"
)

var (
	ErrCurrencyMismatch = errors.New("currency mismatch")
	ErrNegativeAmount   = errors.New("amount cannot be negative")
)

// Money is an immutable value object representing a monetary amount.
type Money struct {
	amount   int64  // stored in smallest currency unit (cents)
	currency string // ISO 4217
}

func NewMoney(amount int64, currency string) (Money, error) {
	if amount < 0 {
		return Money{}, ErrNegativeAmount
	}
	if currency == "" {
		currency = "USD"
	}
	return Money{amount: amount, currency: currency}, nil
}

func MustNewMoney(amount int64, currency string) Money {
	m, err := NewMoney(amount, currency)
	if err != nil {
		panic(err)
	}
	return m
}

func ZeroMoney(currency string) Money {
	return Money{amount: 0, currency: currency}
}

func (m Money) Amount() int64    { return m.amount }
func (m Money) Currency() string { return m.currency }

func (m Money) Add(other Money) (Money, error) {
	if m.currency != other.currency {
		return Money{}, ErrCurrencyMismatch
	}
	return Money{amount: m.amount + other.amount, currency: m.currency}, nil
}

func (m Money) Subtract(other Money) (Money, error) {
	if m.currency != other.currency {
		return Money{}, ErrCurrencyMismatch
	}
	result := m.amount - other.amount
	if result < 0 {
		return Money{}, ErrNegativeAmount
	}
	return Money{amount: result, currency: m.currency}, nil
}

func (m Money) Multiply(factor float64) Money {
	return Money{
		amount:   int64(math.Round(float64(m.amount) * factor)),
		currency: m.currency,
	}
}

func (m Money) IsZero() bool {
	return m.amount == 0
}

func (m Money) GreaterThan(other Money) bool {
	return m.currency == other.currency && m.amount > other.amount
}

func (m Money) String() string {
	whole := m.amount / 100
	frac := m.amount % 100
	return fmt.Sprintf("%s %d.%02d", m.currency, whole, frac)
}

package exchangerate

import (
	"time"

	"github.com/google/uuid"
)

// ExchangeRate stores historical FX rates for currency conversion.
// Typically USD as base with local currencies.
type ExchangeRate struct {
	id           uuid.UUID
	baseCurrency string  // e.g. "USD"
	quoteCurrency string // e.g. "ARS", "BRL", "COP", "MXN"
	rate         float64 // 1 base = rate quote (e.g. 1 USD = 950 ARS)
	inverseRate  float64 // 1 quote = inverseRate base
	source       string  // "manual", "api", "import"
	effectiveDate time.Time
	createdAt    time.Time
}

func NewExchangeRate(
	baseCurrency, quoteCurrency string, rate float64,
	source string, effectiveDate time.Time,
) (*ExchangeRate, error) {
	if baseCurrency == "" || quoteCurrency == "" {
		return nil, ErrInvalidPair
	}
	if rate <= 0 {
		return nil, ErrInvalidRate
	}

	return &ExchangeRate{
		id:            uuid.New(),
		baseCurrency:  baseCurrency,
		quoteCurrency: quoteCurrency,
		rate:          rate,
		inverseRate:   1.0 / rate,
		source:        source,
		effectiveDate: effectiveDate,
		createdAt:     time.Now().UTC(),
	}, nil
}

func HydrateExchangeRate(
	id uuid.UUID, baseCurrency, quoteCurrency string,
	rate, inverseRate float64, source string,
	effectiveDate, createdAt time.Time,
) *ExchangeRate {
	return &ExchangeRate{
		id: id, baseCurrency: baseCurrency, quoteCurrency: quoteCurrency,
		rate: rate, inverseRate: inverseRate, source: source,
		effectiveDate: effectiveDate, createdAt: createdAt,
	}
}

func (r *ExchangeRate) ID() uuid.UUID          { return r.id }
func (r *ExchangeRate) BaseCurrency() string    { return r.baseCurrency }
func (r *ExchangeRate) QuoteCurrency() string   { return r.quoteCurrency }
func (r *ExchangeRate) Rate() float64           { return r.rate }
func (r *ExchangeRate) InverseRate() float64    { return r.inverseRate }
func (r *ExchangeRate) Source() string          { return r.source }
func (r *ExchangeRate) EffectiveDate() time.Time { return r.effectiveDate }
func (r *ExchangeRate) CreatedAt() time.Time    { return r.createdAt }

// Convert converts an amount from base to quote currency.
func (r *ExchangeRate) Convert(amountCents int64) int64 {
	return int64(float64(amountCents) * r.rate)
}

// ConvertInverse converts an amount from quote to base currency.
func (r *ExchangeRate) ConvertInverse(amountCents int64) int64 {
	return int64(float64(amountCents) * r.inverseRate)
}

package response

import (
	"time"

	"github.com/diogenes/costforensics/backend/internal/domain/exchangerate"
)

type ExchangeRate struct {
	ID            string    `json:"id"`
	BaseCurrency  string    `json:"base_currency"`
	QuoteCurrency string    `json:"quote_currency"`
	Rate          float64   `json:"rate"`
	InverseRate   float64   `json:"inverse_rate"`
	Source        string    `json:"source"`
	EffectiveDate string    `json:"effective_date"`
	CreatedAt     time.Time `json:"created_at"`
}

func ExchangeRateFromDomain(r *exchangerate.ExchangeRate) ExchangeRate {
	return ExchangeRate{
		ID: r.ID().String(), BaseCurrency: r.BaseCurrency(),
		QuoteCurrency: r.QuoteCurrency(), Rate: r.Rate(),
		InverseRate: r.InverseRate(), Source: r.Source(),
		EffectiveDate: r.EffectiveDate().Format("2006-01-02"),
		CreatedAt: r.CreatedAt(),
	}
}

package mapper

import (
	"github.com/diogenes/costforensics/backend/internal/adapter/postgres/model"
	"github.com/diogenes/costforensics/backend/internal/domain/exchangerate"
)

func ExchangeRateToModel(r *exchangerate.ExchangeRate) *model.ExchangeRate {
	return &model.ExchangeRate{
		ID:            r.ID(),
		BaseCurrency:  r.BaseCurrency(),
		QuoteCurrency: r.QuoteCurrency(),
		Rate:          r.Rate(),
		InverseRate:   r.InverseRate(),
		Source:        r.Source(),
		EffectiveDate: r.EffectiveDate(),
		CreatedAt:     r.CreatedAt(),
	}
}

func ExchangeRateToDomain(m *model.ExchangeRate) *exchangerate.ExchangeRate {
	return exchangerate.HydrateExchangeRate(
		m.ID, m.BaseCurrency, m.QuoteCurrency,
		m.Rate, m.InverseRate, m.Source,
		m.EffectiveDate, m.CreatedAt,
	)
}

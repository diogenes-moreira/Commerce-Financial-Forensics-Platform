package mapper

import (
	"encoding/json"

	"github.com/diogenes/costforensics/backend/internal/adapter/postgres/model"
	"github.com/diogenes/costforensics/backend/internal/domain/payment"
	"github.com/diogenes/costforensics/backend/internal/domain/shared"
)

func PaymentToModel(p *payment.Payment) *model.Payment {
	metaJSON, _ := json.Marshal(p.Metadata())
	return &model.Payment{
		ID:             p.ID(),
		OrderID:        p.OrderID(),
		Direction:      string(p.Direction()),
		CounterpartyID: p.CounterpartyID(),
		AmountCents:    p.AmountCents().Amount(),
		Currency:       p.AmountCents().Currency(),
		Method:         p.Method(),
		ExternalRef:    p.ExternalRef(),
		Status:         string(p.Status()),
		ProcessedAt:    p.ProcessedAt(),
		FailureReason:  p.FailureReason(),
		Metadata:       metaJSON,
		CreatedAt:      p.CreatedAt(),
		UpdatedAt:      p.UpdatedAt(),
	}
}

func PaymentToDomain(m *model.Payment) *payment.Payment {
	amount := shared.MustNewMoney(m.AmountCents, m.Currency)
	var meta map[string]string
	_ = json.Unmarshal(m.Metadata, &meta)
	if meta == nil {
		meta = make(map[string]string)
	}
	return payment.HydratePayment(
		m.ID, m.OrderID, payment.Direction(m.Direction), m.CounterpartyID,
		amount, m.Method, m.ExternalRef, payment.Status(m.Status),
		m.ProcessedAt, m.FailureReason,
		meta, m.CreatedAt, m.UpdatedAt,
	)
}

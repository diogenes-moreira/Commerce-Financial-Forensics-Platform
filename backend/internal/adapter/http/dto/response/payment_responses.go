package response

import (
	"time"

	"github.com/diogenes/costforensics/backend/internal/domain/payment"
)

type Payment struct {
	ID             string            `json:"id"`
	OrderID        string            `json:"order_id"`
	Direction      string            `json:"direction"`
	CounterpartyID string            `json:"counterparty_id"`
	AmountCents    int64             `json:"amount_cents"`
	Currency       string            `json:"currency"`
	Method         string            `json:"method"`
	ExternalRef    string            `json:"external_ref"`
	Status         string            `json:"status"`
	ProcessedAt    *time.Time        `json:"processed_at"`
	FailureReason  string            `json:"failure_reason,omitempty"`
	Metadata       map[string]string `json:"metadata"`
	CreatedAt      time.Time         `json:"created_at"`
	UpdatedAt      time.Time         `json:"updated_at"`
}

func PaymentFromDomain(p *payment.Payment) Payment {
	return Payment{
		ID: p.ID().String(), OrderID: p.OrderID().String(),
		Direction: string(p.Direction()), CounterpartyID: p.CounterpartyID().String(),
		AmountCents: p.AmountCents().Amount(), Currency: p.AmountCents().Currency(),
		Method: p.Method(), ExternalRef: p.ExternalRef(),
		Status: string(p.Status()), ProcessedAt: p.ProcessedAt(),
		FailureReason: p.FailureReason(), Metadata: p.Metadata(),
		CreatedAt: p.CreatedAt(), UpdatedAt: p.UpdatedAt(),
	}
}

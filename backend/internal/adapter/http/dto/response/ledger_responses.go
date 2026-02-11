package response

import (
	"time"

	"github.com/diogenes/costforensics/backend/internal/domain/ledger"
)

type LedgerEntry struct {
	ID            string    `json:"id"`
	OrderID       *string   `json:"order_id"`
	AccountCode   string    `json:"account_code"`
	Side          string    `json:"side"`
	AmountCents   int64     `json:"amount_cents"`
	Currency      string    `json:"currency"`
	Description   string    `json:"description"`
	ReferenceType string    `json:"reference_type"`
	ReferenceID   string    `json:"reference_id"`
	EffectiveDate string    `json:"effective_date"`
	CreatedAt     time.Time `json:"created_at"`
}

func LedgerEntryFromDomain(e *ledger.LedgerEntry) LedgerEntry {
	var orderID *string
	if e.OrderID() != nil {
		s := e.OrderID().String()
		orderID = &s
	}
	return LedgerEntry{
		ID: e.ID().String(), OrderID: orderID,
		AccountCode: e.AccountCode(), Side: string(e.Side()),
		AmountCents: e.AmountCents().Amount(), Currency: e.AmountCents().Currency(),
		Description: e.Description(), ReferenceType: e.ReferenceType(),
		ReferenceID: e.ReferenceID().String(),
		EffectiveDate: e.EffectiveDate().Format("2006-01-02"),
		CreatedAt: e.CreatedAt(),
	}
}

type BalanceResponse struct {
	AccountCode string `json:"account_code"`
	BalanceCents int64 `json:"balance_cents"`
	StartDate   string `json:"start_date,omitempty"`
	EndDate     string `json:"end_date,omitempty"`
}

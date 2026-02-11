package mapper

import (
	"github.com/diogenes/costforensics/backend/internal/adapter/postgres/model"
	"github.com/diogenes/costforensics/backend/internal/domain/ledger"
	"github.com/diogenes/costforensics/backend/internal/domain/shared"
)

func LedgerEntryToModel(e *ledger.LedgerEntry) *model.LedgerEntry {
	return &model.LedgerEntry{
		ID:            e.ID(),
		OrderID:       e.OrderID(),
		AccountCode:   e.AccountCode(),
		Side:          string(e.Side()),
		AmountCents:   e.AmountCents().Amount(),
		Currency:      e.AmountCents().Currency(),
		Description:   e.Description(),
		ReferenceType: e.ReferenceType(),
		ReferenceID:   e.ReferenceID(),
		EffectiveDate: e.EffectiveDate(),
		CreatedAt:     e.CreatedAt(),
	}
}

func LedgerEntryToDomain(m *model.LedgerEntry) *ledger.LedgerEntry {
	amount := shared.MustNewMoney(m.AmountCents, m.Currency)
	return ledger.HydrateLedgerEntry(
		m.ID, m.OrderID, m.AccountCode, ledger.Side(m.Side),
		amount, m.Description, m.ReferenceType, m.ReferenceID,
		m.EffectiveDate, m.CreatedAt,
	)
}

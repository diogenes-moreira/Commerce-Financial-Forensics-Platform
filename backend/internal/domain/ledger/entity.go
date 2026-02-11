package ledger

import (
	"time"

	"github.com/diogenes/costforensics/backend/internal/domain/shared"
	"github.com/google/uuid"
)

type Side string

const (
	SideDebit  Side = "debit"
	SideCredit Side = "credit"
)

type LedgerEntry struct {
	id            uuid.UUID
	orderID       *uuid.UUID
	accountCode   string
	side          Side
	amountCents   shared.Money
	description   string
	referenceType string
	referenceID   uuid.UUID
	effectiveDate time.Time
	createdAt     time.Time
}

func NewLedgerEntry(
	orderID *uuid.UUID, accountCode string, side Side,
	amount shared.Money, description, refType string,
	refID uuid.UUID, effectiveDate time.Time,
) (*LedgerEntry, error) {
	if accountCode == "" {
		return nil, ErrInvalidAccountCode
	}
	if side != SideDebit && side != SideCredit {
		return nil, ErrInvalidSide
	}

	return &LedgerEntry{
		id:            uuid.New(),
		orderID:       orderID,
		accountCode:   accountCode,
		side:          side,
		amountCents:   amount,
		description:   description,
		referenceType: refType,
		referenceID:   refID,
		effectiveDate: effectiveDate,
		createdAt:     time.Now().UTC(),
	}, nil
}

func HydrateLedgerEntry(
	id uuid.UUID, orderID *uuid.UUID, accountCode string, side Side,
	amount shared.Money, description, refType string, refID uuid.UUID,
	effectiveDate, createdAt time.Time,
) *LedgerEntry {
	return &LedgerEntry{
		id: id, orderID: orderID, accountCode: accountCode, side: side,
		amountCents: amount, description: description,
		referenceType: refType, referenceID: refID,
		effectiveDate: effectiveDate, createdAt: createdAt,
	}
}

func (e *LedgerEntry) ID() uuid.UUID           { return e.id }
func (e *LedgerEntry) OrderID() *uuid.UUID      { return e.orderID }
func (e *LedgerEntry) AccountCode() string      { return e.accountCode }
func (e *LedgerEntry) Side() Side               { return e.side }
func (e *LedgerEntry) AmountCents() shared.Money { return e.amountCents }
func (e *LedgerEntry) Description() string      { return e.description }
func (e *LedgerEntry) ReferenceType() string    { return e.referenceType }
func (e *LedgerEntry) ReferenceID() uuid.UUID   { return e.referenceID }
func (e *LedgerEntry) EffectiveDate() time.Time { return e.effectiveDate }
func (e *LedgerEntry) CreatedAt() time.Time     { return e.createdAt }

package seller

import (
	"time"

	"github.com/diogenes/costforensics/backend/internal/domain/shared"
	"github.com/google/uuid"
)

type Status string

const (
	StatusActive   Status = "active"
	StatusInactive Status = "inactive"
)

type Seller struct {
	id            uuid.UUID
	externalID    string
	code          string // seller code e.g. "a19999"
	name          string
	email         string
	commissionPct float64
	status        Status
	createdAt     time.Time
	updatedAt     time.Time
}

func NewSeller(externalID, code, name, email string, commissionPct float64) (*Seller, error) {
	if name == "" {
		return nil, ErrNameEmpty
	}
	if commissionPct < 0 || commissionPct > 100 {
		return nil, ErrInvalidCommission
	}

	now := time.Now().UTC()
	return &Seller{
		id:            uuid.New(),
		externalID:    externalID,
		code:          code,
		name:          name,
		email:         email,
		commissionPct: commissionPct,
		status:        StatusActive,
		createdAt:     now,
		updatedAt:     now,
	}, nil
}

func HydrateSeller(
	id uuid.UUID, externalID, code, name, email string,
	commissionPct float64, status Status,
	createdAt, updatedAt time.Time,
) *Seller {
	return &Seller{
		id: id, externalID: externalID, code: code, name: name, email: email,
		commissionPct: commissionPct, status: status,
		createdAt: createdAt, updatedAt: updatedAt,
	}
}

func (s *Seller) ID() uuid.UUID        { return s.id }
func (s *Seller) ExternalID() string    { return s.externalID }
func (s *Seller) Code() string          { return s.code }
func (s *Seller) Name() string          { return s.name }
func (s *Seller) Email() string         { return s.email }
func (s *Seller) CommissionPct() float64 { return s.commissionPct }
func (s *Seller) Status() Status        { return s.status }
func (s *Seller) CreatedAt() time.Time  { return s.createdAt }
func (s *Seller) UpdatedAt() time.Time  { return s.updatedAt }

func (s *Seller) UpdateCommission(pct float64) error {
	if pct < 0 || pct > 100 {
		return ErrInvalidCommission
	}
	s.commissionPct = pct
	s.updatedAt = time.Now().UTC()
	return nil
}

func (s *Seller) CalculateCommission(amount shared.Money) shared.Money {
	return amount.Multiply(s.commissionPct / 100)
}

func (s *Seller) Deactivate() {
	s.status = StatusInactive
	s.updatedAt = time.Now().UTC()
}

func (s *Seller) UpdateDetails(name, email string) error {
	if name == "" {
		return ErrNameEmpty
	}
	s.name = name
	s.email = email
	s.updatedAt = time.Now().UTC()
	return nil
}

package customer

import (
	"time"

	"github.com/google/uuid"
)

type Customer struct {
	id         uuid.UUID
	externalID string
	name       string
	email      string
	segment    string
	metadata   map[string]string
	createdAt  time.Time
	updatedAt  time.Time
}

func NewCustomer(externalID, name, email, segment string) (*Customer, error) {
	if name == "" {
		return nil, ErrNameEmpty
	}

	now := time.Now().UTC()
	return &Customer{
		id:         uuid.New(),
		externalID: externalID,
		name:       name,
		email:      email,
		segment:    segment,
		metadata:   make(map[string]string),
		createdAt:  now,
		updatedAt:  now,
	}, nil
}

func HydrateCustomer(
	id uuid.UUID, externalID, name, email, segment string,
	metadata map[string]string, createdAt, updatedAt time.Time,
) *Customer {
	if metadata == nil {
		metadata = make(map[string]string)
	}
	return &Customer{
		id: id, externalID: externalID, name: name, email: email,
		segment: segment, metadata: metadata,
		createdAt: createdAt, updatedAt: updatedAt,
	}
}

func (c *Customer) ID() uuid.UUID              { return c.id }
func (c *Customer) ExternalID() string          { return c.externalID }
func (c *Customer) Name() string                { return c.name }
func (c *Customer) Email() string               { return c.email }
func (c *Customer) Segment() string             { return c.segment }
func (c *Customer) Metadata() map[string]string { return c.metadata }
func (c *Customer) CreatedAt() time.Time        { return c.createdAt }
func (c *Customer) UpdatedAt() time.Time        { return c.updatedAt }

func (c *Customer) UpdateDetails(name, email, segment string) error {
	if name == "" {
		return ErrNameEmpty
	}
	c.name = name
	c.email = email
	c.segment = segment
	c.updatedAt = time.Now().UTC()
	return nil
}

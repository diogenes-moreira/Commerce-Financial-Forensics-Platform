package forensicevent

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type ForensicEvent struct {
	id            uuid.UUID
	entityType    string
	entityID      uuid.UUID
	eventType     string
	eventTime     time.Time
	effectiveTime time.Time
	payload       map[string]any
	actorID       string
	hashIntegrity string
	createdAt     time.Time
}

func NewForensicEvent(
	entityType string, entityID uuid.UUID, eventType string,
	effectiveTime time.Time, payload map[string]any, actorID string,
) (*ForensicEvent, error) {
	if entityType == "" {
		return nil, ErrEntityTypeEmpty
	}
	if eventType == "" {
		return nil, ErrEventTypeEmpty
	}
	if payload == nil {
		payload = make(map[string]any)
	}

	hash := computeHash(payload)
	now := time.Now().UTC()
	return &ForensicEvent{
		id:            uuid.New(),
		entityType:    entityType,
		entityID:      entityID,
		eventType:     eventType,
		eventTime:     now,
		effectiveTime: effectiveTime,
		payload:       payload,
		actorID:       actorID,
		hashIntegrity: hash,
		createdAt:     now,
	}, nil
}

func HydrateForensicEvent(
	id uuid.UUID, entityType string, entityID uuid.UUID, eventType string,
	eventTime, effectiveTime time.Time, payload map[string]any,
	actorID, hashIntegrity string, createdAt time.Time,
) *ForensicEvent {
	if payload == nil {
		payload = make(map[string]any)
	}
	return &ForensicEvent{
		id: id, entityType: entityType, entityID: entityID,
		eventType: eventType, eventTime: eventTime, effectiveTime: effectiveTime,
		payload: payload, actorID: actorID, hashIntegrity: hashIntegrity,
		createdAt: createdAt,
	}
}

func (e *ForensicEvent) ID() uuid.UUID            { return e.id }
func (e *ForensicEvent) EntityType() string        { return e.entityType }
func (e *ForensicEvent) EntityID() uuid.UUID       { return e.entityID }
func (e *ForensicEvent) EventType() string         { return e.eventType }
func (e *ForensicEvent) EventTime() time.Time      { return e.eventTime }
func (e *ForensicEvent) EffectiveTime() time.Time  { return e.effectiveTime }
func (e *ForensicEvent) Payload() map[string]any   { return e.payload }
func (e *ForensicEvent) ActorID() string           { return e.actorID }
func (e *ForensicEvent) HashIntegrity() string     { return e.hashIntegrity }
func (e *ForensicEvent) CreatedAt() time.Time      { return e.createdAt }

func computeHash(payload map[string]any) string {
	data, _ := json.Marshal(payload)
	h := sha256.Sum256(data)
	return fmt.Sprintf("%x", h)
}

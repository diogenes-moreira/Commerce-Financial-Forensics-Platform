package response

import (
	"time"

	"github.com/diogenes/costforensics/backend/internal/domain/forensicevent"
)

type ForensicEvent struct {
	ID            string         `json:"id"`
	EntityType    string         `json:"entity_type"`
	EntityID      string         `json:"entity_id"`
	EventType     string         `json:"event_type"`
	EventTime     time.Time      `json:"event_time"`
	EffectiveTime time.Time      `json:"effective_time"`
	Payload       map[string]any `json:"payload"`
	ActorID       string         `json:"actor_id"`
	HashIntegrity string         `json:"hash_integrity"`
	CreatedAt     time.Time      `json:"created_at"`
}

func ForensicEventFromDomain(e *forensicevent.ForensicEvent) ForensicEvent {
	return ForensicEvent{
		ID: e.ID().String(), EntityType: e.EntityType(),
		EntityID: e.EntityID().String(), EventType: e.EventType(),
		EventTime: e.EventTime(), EffectiveTime: e.EffectiveTime(),
		Payload: e.Payload(), ActorID: e.ActorID(),
		HashIntegrity: e.HashIntegrity(), CreatedAt: e.CreatedAt(),
	}
}

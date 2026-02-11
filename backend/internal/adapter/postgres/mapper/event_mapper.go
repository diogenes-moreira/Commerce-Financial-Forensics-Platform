package mapper

import (
	"encoding/json"

	"github.com/diogenes/costforensics/backend/internal/adapter/postgres/model"
	"github.com/diogenes/costforensics/backend/internal/domain/forensicevent"
)

func ForensicEventToModel(e *forensicevent.ForensicEvent) *model.ForensicEvent {
	payloadJSON, _ := json.Marshal(e.Payload())
	return &model.ForensicEvent{
		ID:            e.ID(),
		EntityType:    e.EntityType(),
		EntityID:      e.EntityID(),
		EventType:     e.EventType(),
		EventTime:     e.EventTime(),
		EffectiveTime: e.EffectiveTime(),
		Payload:       payloadJSON,
		ActorID:       e.ActorID(),
		HashIntegrity: e.HashIntegrity(),
		CreatedAt:     e.CreatedAt(),
	}
}

func ForensicEventToDomain(m *model.ForensicEvent) *forensicevent.ForensicEvent {
	var payload map[string]any
	_ = json.Unmarshal(m.Payload, &payload)
	if payload == nil {
		payload = make(map[string]any)
	}
	return forensicevent.HydrateForensicEvent(
		m.ID, m.EntityType, m.EntityID, m.EventType,
		m.EventTime, m.EffectiveTime, payload,
		m.ActorID, m.HashIntegrity, m.CreatedAt,
	)
}

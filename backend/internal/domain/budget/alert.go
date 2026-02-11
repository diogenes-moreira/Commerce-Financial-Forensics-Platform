package budget

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Alert struct {
	id           uuid.UUID
	budgetID     uuid.UUID
	thresholdPct int
	actualPct    int
	message      string
	triggeredAt  time.Time
}

func newAlert(budgetID uuid.UUID, thresholdPct, actualPct int) *Alert {
	return &Alert{
		id:           uuid.New(),
		budgetID:     budgetID,
		thresholdPct: thresholdPct,
		actualPct:    actualPct,
		message:      fmt.Sprintf("Budget usage at %d%% (threshold: %d%%)", actualPct, thresholdPct),
		triggeredAt:  time.Now().UTC(),
	}
}

func HydrateAlert(id, budgetID uuid.UUID, thresholdPct, actualPct int, message string, triggeredAt time.Time) *Alert {
	return &Alert{
		id: id, budgetID: budgetID, thresholdPct: thresholdPct,
		actualPct: actualPct, message: message, triggeredAt: triggeredAt,
	}
}

func (a *Alert) ID() uuid.UUID       { return a.id }
func (a *Alert) BudgetID() uuid.UUID  { return a.budgetID }
func (a *Alert) ThresholdPct() int    { return a.thresholdPct }
func (a *Alert) ActualPct() int       { return a.actualPct }
func (a *Alert) Message() string      { return a.message }
func (a *Alert) TriggeredAt() time.Time { return a.triggeredAt }

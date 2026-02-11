package anomaly

import (
	"time"

	"github.com/diogenes/costforensics/backend/internal/domain/shared"
	"github.com/google/uuid"
)

type Severity string

const (
	SeverityLow      Severity = "low"
	SeverityMedium   Severity = "medium"
	SeverityHigh     Severity = "high"
	SeverityCritical Severity = "critical"
)

type Status string

const (
	StatusOpen     Status = "open"
	StatusResolved Status = "resolved"
)

type CostAnomaly struct {
	id             uuid.UUID
	cloudAccountID uuid.UUID
	service        string
	expected       shared.Money
	actual         shared.Money
	deviationPct   float64
	severity       Severity
	status         Status
	detectedAt     time.Time
	resolvedAt     *time.Time
}

func NewCostAnomaly(cloudAccountID uuid.UUID, service string, expected, actual shared.Money, deviationPct float64) (*CostAnomaly, error) {
	now := time.Now().UTC()
	a := &CostAnomaly{
		id:             uuid.New(),
		cloudAccountID: cloudAccountID,
		service:        service,
		expected:       expected,
		actual:         actual,
		deviationPct:   deviationPct,
		status:         StatusOpen,
		detectedAt:     now,
	}
	a.severity = a.classifySeverity()
	return a, nil
}

func HydrateAnomaly(
	id, cloudAccountID uuid.UUID, service string,
	expected, actual shared.Money, deviationPct float64,
	severity Severity, status Status,
	detectedAt time.Time, resolvedAt *time.Time,
) *CostAnomaly {
	return &CostAnomaly{
		id: id, cloudAccountID: cloudAccountID, service: service,
		expected: expected, actual: actual, deviationPct: deviationPct,
		severity: severity, status: status,
		detectedAt: detectedAt, resolvedAt: resolvedAt,
	}
}

func (a *CostAnomaly) ID() uuid.UUID             { return a.id }
func (a *CostAnomaly) CloudAccountID() uuid.UUID  { return a.cloudAccountID }
func (a *CostAnomaly) Service() string             { return a.service }
func (a *CostAnomaly) Expected() shared.Money      { return a.expected }
func (a *CostAnomaly) Actual() shared.Money        { return a.actual }
func (a *CostAnomaly) DeviationPct() float64       { return a.deviationPct }
func (a *CostAnomaly) Severity() Severity          { return a.severity }
func (a *CostAnomaly) Status() Status              { return a.status }
func (a *CostAnomaly) DetectedAt() time.Time       { return a.detectedAt }
func (a *CostAnomaly) ResolvedAt() *time.Time      { return a.resolvedAt }

// classifySeverity self-classifies based on deviation percentage.
func (a *CostAnomaly) classifySeverity() Severity {
	switch {
	case a.deviationPct >= 200:
		return SeverityCritical
	case a.deviationPct >= 100:
		return SeverityHigh
	case a.deviationPct >= 50:
		return SeverityMedium
	default:
		return SeverityLow
	}
}

func (a *CostAnomaly) Resolve() error {
	if a.status == StatusResolved {
		return ErrAlreadyResolved
	}
	now := time.Now().UTC()
	a.status = StatusResolved
	a.resolvedAt = &now
	return nil
}

package anomaly

import (
	"github.com/diogenes/costforensics/backend/internal/domain/shared"
	"github.com/google/uuid"
)

// Detector detects cost anomalies by comparing actual vs expected.
type Detector struct{}

func NewDetector() *Detector {
	return &Detector{}
}

// DetectAnomaly compares expected and actual costs. Returns anomaly if deviation
// exceeds the given threshold percentage.
func (d *Detector) DetectAnomaly(cloudAccountID uuid.UUID, service string, expected, actual shared.Money, thresholdPct float64) *CostAnomaly {
	if expected.IsZero() {
		return nil
	}

	deviationPct := float64(actual.Amount()-expected.Amount()) / float64(expected.Amount()) * 100

	if deviationPct < thresholdPct {
		return nil
	}

	a, _ := NewCostAnomaly(cloudAccountID, service, expected, actual, deviationPct)
	return a
}

package costreport

import (
	"time"

	"github.com/diogenes/costforensics/backend/internal/domain/shared"
	"github.com/google/uuid"
)

type ReportType string

const (
	ReportTypeDaily    ReportType = "daily"
	ReportTypeWeekly   ReportType = "weekly"
	ReportTypeMonthly  ReportType = "monthly"
	ReportTypeCustom   ReportType = "custom"
)

type CostReport struct {
	id          uuid.UUID
	name        string
	reportType  ReportType
	periodStart time.Time
	periodEnd   time.Time
	total       shared.Money
	breakdown   map[string]int64 // service -> amount in cents
	generatedAt time.Time
}

func NewCostReport(name string, reportType ReportType, periodStart, periodEnd time.Time, currency string) (*CostReport, error) {
	if name == "" {
		return nil, ErrReportNameEmpty
	}
	if !periodStart.Before(periodEnd) {
		return nil, ErrInvalidPeriod
	}

	return &CostReport{
		id:          uuid.New(),
		name:        name,
		reportType:  reportType,
		periodStart: periodStart,
		periodEnd:   periodEnd,
		total:       shared.ZeroMoney(currency),
		breakdown:   make(map[string]int64),
		generatedAt: time.Now().UTC(),
	}, nil
}

func HydrateCostReport(
	id uuid.UUID, name string, reportType ReportType,
	periodStart, periodEnd time.Time, total shared.Money,
	breakdown map[string]int64, generatedAt time.Time,
) *CostReport {
	return &CostReport{
		id: id, name: name, reportType: reportType,
		periodStart: periodStart, periodEnd: periodEnd,
		total: total, breakdown: breakdown, generatedAt: generatedAt,
	}
}

func (r *CostReport) ID() uuid.UUID                { return r.id }
func (r *CostReport) Name() string                  { return r.name }
func (r *CostReport) ReportType() ReportType        { return r.reportType }
func (r *CostReport) PeriodStart() time.Time        { return r.periodStart }
func (r *CostReport) PeriodEnd() time.Time          { return r.periodEnd }
func (r *CostReport) Total() shared.Money           { return r.total }
func (r *CostReport) Breakdown() map[string]int64   { return r.breakdown }
func (r *CostReport) GeneratedAt() time.Time        { return r.generatedAt }

// AddServiceCost aggregates cost for a service.
func (r *CostReport) AddServiceCost(service string, amount shared.Money) {
	r.breakdown[service] += amount.Amount()
	r.total, _ = r.total.Add(amount)
}

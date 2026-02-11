package mapper

import (
	"encoding/json"

	"github.com/diogenes/costforensics/backend/internal/adapter/postgres/model"
	"github.com/diogenes/costforensics/backend/internal/domain/anomaly"
	"github.com/diogenes/costforensics/backend/internal/domain/budget"
	"github.com/diogenes/costforensics/backend/internal/domain/cloudaccount"
	"github.com/diogenes/costforensics/backend/internal/domain/costreport"
	"github.com/diogenes/costforensics/backend/internal/domain/costrecord"
	"github.com/diogenes/costforensics/backend/internal/domain/shared"
)

// CloudAccount

func CloudAccountToModel(a *cloudaccount.CloudAccount) *model.CloudAccount {
	return &model.CloudAccount{
		ID:           a.ID(),
		Provider:     string(a.Provider()),
		Name:         a.Name(),
		ExternalID:   a.ExternalID(),
		Status:       string(a.Status()),
		SyncStatus:   string(a.SyncStatus()),
		LastSyncedAt: a.LastSyncedAt(),
		CreatedAt:    a.CreatedAt(),
		UpdatedAt:    a.UpdatedAt(),
	}
}

func CloudAccountToDomain(m *model.CloudAccount) *cloudaccount.CloudAccount {
	return cloudaccount.HydrateCloudAccount(
		m.ID, cloudaccount.Provider(m.Provider), m.Name, m.ExternalID,
		cloudaccount.Status(m.Status), cloudaccount.SyncStatus(m.SyncStatus),
		m.LastSyncedAt, m.CreatedAt, m.UpdatedAt,
	)
}

// CostRecord

func CostRecordToModel(r *costrecord.CostRecord) *model.CostRecord {
	tagsJSON, _ := json.Marshal(r.Tags())
	return &model.CostRecord{
		ID:             r.ID(),
		CloudAccountID: r.CloudAccountID(),
		Service:        r.Service(),
		Category:       string(r.Category()),
		AmountCents:    r.Amount().Amount(),
		Currency:       r.Amount().Currency(),
		UsageDate:      r.UsageDate(),
		Tags:           tagsJSON,
		CreatedAt:      r.CreatedAt(),
	}
}

func CostRecordToDomain(m *model.CostRecord) *costrecord.CostRecord {
	amount := shared.MustNewMoney(m.AmountCents, m.Currency)
	var tags map[string]string
	_ = json.Unmarshal(m.Tags, &tags)
	if tags == nil {
		tags = make(map[string]string)
	}
	return costrecord.HydrateCostRecord(
		m.ID, m.CloudAccountID, m.Service,
		costrecord.Category(m.Category), amount,
		m.UsageDate, tags, m.CreatedAt,
	)
}

// Budget

func BudgetToModel(b *budget.Budget) *model.Budget {
	return &model.Budget{
		ID:               b.ID(),
		Name:             b.Name(),
		AmountCents:      b.Amount().Amount(),
		Currency:         b.Amount().Currency(),
		SpentCents:       b.Spent().Amount(),
		PeriodStart:      b.PeriodStart(),
		PeriodEnd:        b.PeriodEnd(),
		AlertThresholdPct: b.AlertThresholdPct(),
		Status:           string(b.Status()),
		CreatedAt:        b.CreatedAt(),
		UpdatedAt:        b.UpdatedAt(),
	}
}

func BudgetToDomain(m *model.Budget) *budget.Budget {
	amount := shared.MustNewMoney(m.AmountCents, m.Currency)
	spent := shared.MustNewMoney(m.SpentCents, m.Currency)
	return budget.HydrateBudget(
		m.ID, m.Name, amount, spent,
		m.PeriodStart, m.PeriodEnd, m.AlertThresholdPct,
		budget.Status(m.Status), m.CreatedAt, m.UpdatedAt,
	)
}

func BudgetAlertToModel(a *budget.Alert) *model.BudgetAlert {
	return &model.BudgetAlert{
		ID:           a.ID(),
		BudgetID:     a.BudgetID(),
		ThresholdPct: a.ThresholdPct(),
		ActualPct:    a.ActualPct(),
		Message:      a.Message(),
		TriggeredAt:  a.TriggeredAt(),
	}
}

// Anomaly

func AnomalyToModel(a *anomaly.CostAnomaly) *model.Anomaly {
	return &model.Anomaly{
		ID:             a.ID(),
		CloudAccountID: a.CloudAccountID(),
		Service:        a.Service(),
		ExpectedCents:  a.Expected().Amount(),
		ActualCents:    a.Actual().Amount(),
		Currency:       a.Expected().Currency(),
		DeviationPct:   a.DeviationPct(),
		Severity:       string(a.Severity()),
		Status:         string(a.Status()),
		DetectedAt:     a.DetectedAt(),
		ResolvedAt:     a.ResolvedAt(),
	}
}

func AnomalyToDomain(m *model.Anomaly) *anomaly.CostAnomaly {
	expected := shared.MustNewMoney(m.ExpectedCents, m.Currency)
	actual := shared.MustNewMoney(m.ActualCents, m.Currency)
	return anomaly.HydrateAnomaly(
		m.ID, m.CloudAccountID, m.Service,
		expected, actual, m.DeviationPct,
		anomaly.Severity(m.Severity), anomaly.Status(m.Status),
		m.DetectedAt, m.ResolvedAt,
	)
}

// CostReport

func CostReportToModel(r *costreport.CostReport) *model.CostReport {
	breakdownJSON, _ := json.Marshal(r.Breakdown())
	return &model.CostReport{
		ID:          r.ID(),
		Name:        r.Name(),
		ReportType:  string(r.ReportType()),
		PeriodStart: r.PeriodStart(),
		PeriodEnd:   r.PeriodEnd(),
		TotalCents:  r.Total().Amount(),
		Currency:    r.Total().Currency(),
		Breakdown:   breakdownJSON,
		GeneratedAt: r.GeneratedAt(),
	}
}

func CostReportToDomain(m *model.CostReport) *costreport.CostReport {
	total := shared.MustNewMoney(m.TotalCents, m.Currency)
	var breakdown map[string]int64
	_ = json.Unmarshal(m.Breakdown, &breakdown)
	if breakdown == nil {
		breakdown = make(map[string]int64)
	}
	return costreport.HydrateCostReport(
		m.ID, m.Name, costreport.ReportType(m.ReportType),
		m.PeriodStart, m.PeriodEnd, total,
		breakdown, m.GeneratedAt,
	)
}

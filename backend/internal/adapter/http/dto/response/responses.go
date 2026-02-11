package response

import (
	"time"

	"github.com/diogenes/costforensics/backend/internal/domain/anomaly"
	"github.com/diogenes/costforensics/backend/internal/domain/budget"
	"github.com/diogenes/costforensics/backend/internal/domain/cloudaccount"
	"github.com/diogenes/costforensics/backend/internal/domain/costreport"
	"github.com/diogenes/costforensics/backend/internal/domain/costrecord"
	"github.com/diogenes/costforensics/backend/internal/domain/tenant"
)

type Tenant struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	DBName    string    `json:"db_name"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func TenantFromDomain(t *tenant.Tenant) Tenant {
	return Tenant{
		ID: t.ID().String(), Name: t.Name(), Slug: t.Slug(),
		DBName: t.DBName(), Status: string(t.Status()),
		CreatedAt: t.CreatedAt(), UpdatedAt: t.UpdatedAt(),
	}
}

type CloudAccount struct {
	ID           string     `json:"id"`
	Provider     string     `json:"provider"`
	Name         string     `json:"name"`
	ExternalID   string     `json:"external_id"`
	Status       string     `json:"status"`
	SyncStatus   string     `json:"sync_status"`
	LastSyncedAt *time.Time `json:"last_synced_at"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

func CloudAccountFromDomain(a *cloudaccount.CloudAccount) CloudAccount {
	return CloudAccount{
		ID: a.ID().String(), Provider: string(a.Provider()),
		Name: a.Name(), ExternalID: a.ExternalID(),
		Status: string(a.Status()), SyncStatus: string(a.SyncStatus()),
		LastSyncedAt: a.LastSyncedAt(), CreatedAt: a.CreatedAt(), UpdatedAt: a.UpdatedAt(),
	}
}

type CostRecord struct {
	ID             string            `json:"id"`
	CloudAccountID string            `json:"cloud_account_id"`
	Service        string            `json:"service"`
	Category       string            `json:"category"`
	AmountCents    int64             `json:"amount_cents"`
	Currency       string            `json:"currency"`
	UsageDate      string            `json:"usage_date"`
	Tags           map[string]string `json:"tags"`
	CreatedAt      time.Time         `json:"created_at"`
}

func CostRecordFromDomain(r *costrecord.CostRecord) CostRecord {
	return CostRecord{
		ID: r.ID().String(), CloudAccountID: r.CloudAccountID().String(),
		Service: r.Service(), Category: string(r.Category()),
		AmountCents: r.Amount().Amount(), Currency: r.Amount().Currency(),
		UsageDate: r.UsageDate().Format("2006-01-02"), Tags: r.Tags(),
		CreatedAt: r.CreatedAt(),
	}
}

type Budget struct {
	ID               string    `json:"id"`
	Name             string    `json:"name"`
	AmountCents      int64     `json:"amount_cents"`
	Currency         string    `json:"currency"`
	SpentCents       int64     `json:"spent_cents"`
	UsagePct         int       `json:"usage_pct"`
	PeriodStart      string    `json:"period_start"`
	PeriodEnd        string    `json:"period_end"`
	AlertThresholdPct int      `json:"alert_threshold_pct"`
	Status           string    `json:"status"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

func BudgetFromDomain(b *budget.Budget) Budget {
	return Budget{
		ID: b.ID().String(), Name: b.Name(),
		AmountCents: b.Amount().Amount(), Currency: b.Amount().Currency(),
		SpentCents: b.Spent().Amount(), UsagePct: b.UsagePct(),
		PeriodStart: b.PeriodStart().Format("2006-01-02"),
		PeriodEnd: b.PeriodEnd().Format("2006-01-02"),
		AlertThresholdPct: b.AlertThresholdPct(),
		Status: string(b.Status()), CreatedAt: b.CreatedAt(), UpdatedAt: b.UpdatedAt(),
	}
}

type BudgetAlert struct {
	ID           string    `json:"id"`
	BudgetID     string    `json:"budget_id"`
	ThresholdPct int       `json:"threshold_pct"`
	ActualPct    int       `json:"actual_pct"`
	Message      string    `json:"message"`
	TriggeredAt  time.Time `json:"triggered_at"`
}

func BudgetAlertFromDomain(a *budget.Alert) *BudgetAlert {
	if a == nil {
		return nil
	}
	return &BudgetAlert{
		ID: a.ID().String(), BudgetID: a.BudgetID().String(),
		ThresholdPct: a.ThresholdPct(), ActualPct: a.ActualPct(),
		Message: a.Message(), TriggeredAt: a.TriggeredAt(),
	}
}

type Anomaly struct {
	ID             string     `json:"id"`
	CloudAccountID string     `json:"cloud_account_id"`
	Service        string     `json:"service"`
	ExpectedCents  int64      `json:"expected_cents"`
	ActualCents    int64      `json:"actual_cents"`
	Currency       string     `json:"currency"`
	DeviationPct   float64    `json:"deviation_pct"`
	Severity       string     `json:"severity"`
	Status         string     `json:"status"`
	DetectedAt     time.Time  `json:"detected_at"`
	ResolvedAt     *time.Time `json:"resolved_at"`
}

func AnomalyFromDomain(a *anomaly.CostAnomaly) Anomaly {
	return Anomaly{
		ID: a.ID().String(), CloudAccountID: a.CloudAccountID().String(),
		Service: a.Service(), ExpectedCents: a.Expected().Amount(),
		ActualCents: a.Actual().Amount(), Currency: a.Expected().Currency(),
		DeviationPct: a.DeviationPct(), Severity: string(a.Severity()),
		Status: string(a.Status()), DetectedAt: a.DetectedAt(), ResolvedAt: a.ResolvedAt(),
	}
}

type CostReport struct {
	ID          string           `json:"id"`
	Name        string           `json:"name"`
	ReportType  string           `json:"report_type"`
	PeriodStart string           `json:"period_start"`
	PeriodEnd   string           `json:"period_end"`
	TotalCents  int64            `json:"total_cents"`
	Currency    string           `json:"currency"`
	Breakdown   map[string]int64 `json:"breakdown"`
	GeneratedAt time.Time        `json:"generated_at"`
}

func CostReportFromDomain(r *costreport.CostReport) CostReport {
	return CostReport{
		ID: r.ID().String(), Name: r.Name(), ReportType: string(r.ReportType()),
		PeriodStart: r.PeriodStart().Format("2006-01-02"),
		PeriodEnd: r.PeriodEnd().Format("2006-01-02"),
		TotalCents: r.Total().Amount(), Currency: r.Total().Currency(),
		Breakdown: r.Breakdown(), GeneratedAt: r.GeneratedAt(),
	}
}

type Paginated[T any] struct {
	Items      []T   `json:"items"`
	TotalCount int64 `json:"total_count"`
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
}

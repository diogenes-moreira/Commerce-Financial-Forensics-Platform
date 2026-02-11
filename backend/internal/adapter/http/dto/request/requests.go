package request

import "time"

type CreateTenant struct {
	Name string `json:"name" binding:"required"`
	Slug string `json:"slug" binding:"required"`
}

type UpdateTenant struct {
	Name string `json:"name" binding:"required"`
}

type CreateCloudAccount struct {
	Provider   string `json:"provider" binding:"required"`
	Name       string `json:"name" binding:"required"`
	ExternalID string `json:"external_id" binding:"required"`
}

type UpdateCloudAccount struct {
	Name string `json:"name" binding:"required"`
}

type CreateCostRecord struct {
	CloudAccountID string `json:"cloud_account_id" binding:"required"`
	Service        string `json:"service" binding:"required"`
	AmountCents    int64  `json:"amount_cents" binding:"required"`
	Currency       string `json:"currency"`
	UsageDate      string `json:"usage_date" binding:"required"` // YYYY-MM-DD
}

type CreateBudget struct {
	Name             string `json:"name" binding:"required"`
	AmountCents      int64  `json:"amount_cents" binding:"required"`
	Currency         string `json:"currency"`
	PeriodStart      string `json:"period_start" binding:"required"`
	PeriodEnd        string `json:"period_end" binding:"required"`
	AlertThresholdPct int   `json:"alert_threshold_pct"`
}

type UpdateBudget struct {
	Name string `json:"name" binding:"required"`
}

type RecordSpend struct {
	AmountCents int64  `json:"amount_cents" binding:"required"`
	Currency    string `json:"currency"`
}

type GenerateCostReport struct {
	Name        string `json:"name" binding:"required"`
	ReportType  string `json:"report_type" binding:"required"`
	PeriodStart string `json:"period_start" binding:"required"`
	PeriodEnd   string `json:"period_end" binding:"required"`
	Currency    string `json:"currency"`
}

type ListCostRecords struct {
	CloudAccountID string `form:"cloud_account_id"`
	Service        string `form:"service"`
	Category       string `form:"category"`
	StartDate      string `form:"start_date"`
	EndDate        string `form:"end_date"`
	Page           int    `form:"page"`
	PageSize       int    `form:"page_size"`
}

func ParseDate(s string) (time.Time, error) {
	return time.Parse("2006-01-02", s)
}

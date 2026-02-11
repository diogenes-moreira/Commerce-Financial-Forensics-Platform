package request

type CreateExchangeRate struct {
	BaseCurrency  string  `json:"base_currency" binding:"required"`
	QuoteCurrency string  `json:"quote_currency" binding:"required"`
	Rate          float64 `json:"rate" binding:"required"`
	Source        string  `json:"source"`
	EffectiveDate string  `json:"effective_date" binding:"required"`
}

type ListExchangeRates struct {
	BaseCurrency  string `form:"base_currency"`
	QuoteCurrency string `form:"quote_currency"`
	StartDate     string `form:"start_date"`
	EndDate       string `form:"end_date"`
	Page          int    `form:"page"`
	PageSize      int    `form:"page_size"`
}

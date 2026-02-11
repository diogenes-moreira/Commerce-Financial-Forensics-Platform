package request

// ListMargins is the request DTO for listing margin breakdowns within a date range.
type ListMargins struct {
	StartDate string `form:"start_date" binding:"required"`
	EndDate   string `form:"end_date" binding:"required"`
	SellerID  string `form:"seller_id"`
	Page      int    `form:"page"`
	PageSize  int    `form:"page_size"`
}

// ListDrifts is the request DTO for listing drift results within a date range.
type ListDrifts struct {
	StartDate string `form:"start_date" binding:"required"`
	EndDate   string `form:"end_date" binding:"required"`
}

// GetPL is the request DTO for generating a P&L report.
type GetPL struct {
	StartDate   string `form:"start_date" binding:"required"`
	EndDate     string `form:"end_date" binding:"required"`
	Granularity string `form:"granularity"`
	Currency    string `form:"currency"`
}

// GetCohorts is the request DTO for generating a retention cohort report.
type GetCohorts struct {
	StartDate   string `form:"start_date" binding:"required"`
	EndDate     string `form:"end_date" binding:"required"`
	Granularity string `form:"granularity"`
}

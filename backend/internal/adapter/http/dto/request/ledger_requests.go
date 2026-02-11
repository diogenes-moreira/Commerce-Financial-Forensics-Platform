package request

type ListLedgerEntries struct {
	OrderID     string `form:"order_id"`
	AccountCode string `form:"account_code"`
	StartDate   string `form:"start_date"`
	EndDate     string `form:"end_date"`
	Page        int    `form:"page"`
	PageSize    int    `form:"page_size"`
}

type GetBalance struct {
	AccountCode string `form:"account_code" binding:"required"`
	StartDate   string `form:"start_date"`
	EndDate     string `form:"end_date"`
}

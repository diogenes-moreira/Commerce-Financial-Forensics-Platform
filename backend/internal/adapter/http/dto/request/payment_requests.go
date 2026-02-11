package request

type CreatePayment struct {
	OrderID        string `json:"order_id" binding:"required"`
	Direction      string `json:"direction" binding:"required"`
	CounterpartyID string `json:"counterparty_id" binding:"required"`
	AmountCents    int64  `json:"amount_cents" binding:"required"`
	Currency       string `json:"currency"`
	Method         string `json:"method" binding:"required"`
	ExternalRef    string `json:"external_ref"`
}

type ListPayments struct {
	OrderID        string `form:"order_id"`
	Direction      string `form:"direction"`
	CounterpartyID string `form:"counterparty_id"`
	Status         string `form:"status"`
	Page           int    `form:"page"`
	PageSize       int    `form:"page_size"`
}

type FailPayment struct {
	Reason string `json:"reason" binding:"required"`
}

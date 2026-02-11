package request

type CreateProduct struct {
	SKU            string `json:"sku" binding:"required"`
	EAN            string `json:"ean"`
	UPC            string `json:"upc"`
	Name           string `json:"name" binding:"required"`
	Description    string `json:"description"`
	Category       string `json:"category"`
	UnitCostCents  int64  `json:"unit_cost_cents" binding:"required"`
	UnitPriceCents int64  `json:"unit_price_cents" binding:"required"`
	Currency       string `json:"currency"`
}

type UpdateProduct struct {
	Name           string `json:"name" binding:"required"`
	Description    string `json:"description"`
	Category       string `json:"category"`
	UnitPriceCents *int64 `json:"unit_price_cents"`
	Currency       string `json:"currency"`
}

type ListProducts struct {
	Category string `form:"category"`
	Status   string `form:"status"`
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
}

type CreateSeller struct {
	ExternalID    string  `json:"external_id"`
	Code          string  `json:"code"`
	Name          string  `json:"name" binding:"required"`
	Email         string  `json:"email"`
	CommissionPct float64 `json:"commission_pct"`
}

type UpdateSeller struct {
	Name          string   `json:"name" binding:"required"`
	Email         string   `json:"email"`
	CommissionPct *float64 `json:"commission_pct"`
}

type ListSellers struct {
	Status   string `form:"status"`
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
}

type CreateCustomer struct {
	ExternalID string `json:"external_id"`
	Name       string `json:"name" binding:"required"`
	Email      string `json:"email"`
	Segment    string `json:"segment"`
}

type UpdateCustomer struct {
	Name    string `json:"name" binding:"required"`
	Email   string `json:"email"`
	Segment string `json:"segment"`
}

type ListCustomers struct {
	Segment  string `form:"segment"`
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
}

type CreateOrder struct {
	ExternalID string `json:"external_id" binding:"required"`
	SellerID   string `json:"seller_id" binding:"required"`
	CustomerID string `json:"customer_id" binding:"required"`
	Currency   string `json:"currency"`
}

type AddOrderItem struct {
	ProductID      string `json:"product_id" binding:"required"`
	SKU            string `json:"sku" binding:"required"`
	ProductName    string `json:"product_name" binding:"required"`
	Quantity       int    `json:"quantity" binding:"required"`
	UnitPriceCents int64  `json:"unit_price_cents" binding:"required"`
	UnitCostCents  int64  `json:"unit_cost_cents" binding:"required"`
	Currency       string `json:"currency"`
}

type ListOrders struct {
	SellerID   string `form:"seller_id"`
	CustomerID string `form:"customer_id"`
	Status     string `form:"status"`
	StartDate  string `form:"start_date"`
	EndDate    string `form:"end_date"`
	Page       int    `form:"page"`
	PageSize   int    `form:"page_size"`
}

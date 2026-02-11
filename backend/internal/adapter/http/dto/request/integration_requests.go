package request

type CreateIntegration struct {
	Platform string            `json:"platform" binding:"required"`
	Name     string            `json:"name" binding:"required"`
	Config   map[string]string `json:"config"`
}

type UpdateIntegration struct {
	Name   string            `json:"name"`
	Config map[string]string `json:"config"`
}

type ListIntegrations struct {
	Platform string `form:"platform"`
	Status   string `form:"status"`
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
}

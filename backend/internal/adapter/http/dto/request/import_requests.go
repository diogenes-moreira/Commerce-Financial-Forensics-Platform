package request

type ListImports struct {
	Status     string `form:"status"`
	EntityType string `form:"entity_type"`
	Page       int    `form:"page"`
	PageSize   int    `form:"page_size"`
}

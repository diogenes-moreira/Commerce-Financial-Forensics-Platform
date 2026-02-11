package shared

const (
	DefaultPageSize = 20
	MaxPageSize     = 100
)

// Pagination is a value object for paginated queries.
type Pagination struct {
	page     int
	pageSize int
}

func NewPagination(page, pageSize int) Pagination {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = DefaultPageSize
	}
	if pageSize > MaxPageSize {
		pageSize = MaxPageSize
	}
	return Pagination{page: page, pageSize: pageSize}
}

func (p Pagination) Page() int     { return p.page }
func (p Pagination) PageSize() int { return p.pageSize }
func (p Pagination) Offset() int   { return (p.page - 1) * p.pageSize }

type PagedResult[T any] struct {
	Items      []T
	TotalCount int64
	Page       int
	PageSize   int
}

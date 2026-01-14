package repository

type Pagination struct {
	TotalCount int64 `json:"total_count"`
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
}

type PaginatedEntities[T any] struct {
	Entities   []T        `json:"entities"`
	Pagination Pagination `json:"pagination"`
}

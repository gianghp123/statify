package core

import (
	"github.com/gianghp/statify/internal/core/repository"
)

type ApiResponse[T any] struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    *T     `json:"data"`
}

type PaginatedApiResponse[T any] struct {
	Code       int                    `json:"code"`
	Message    string                 `json:"message"`
	Data       []*T                   `json:"data"`
	Pagination *repository.Pagination `json:"pagination"`
}

func NewApiResponse[T any](code int, message string, data *T) *ApiResponse[T] {
	return &ApiResponse[T]{
		Code:    code,
		Message: message,
		Data:    data,
	}
}

func NewPaginatedApiResponse[T any](
	code int,
	message string,
	data []*T,
	pagination *repository.Pagination,
) *PaginatedApiResponse[T] {
	return &PaginatedApiResponse[T]{
		Code:       code,
		Message:    message,
		Data:       data,
		Pagination: pagination,
	}
}

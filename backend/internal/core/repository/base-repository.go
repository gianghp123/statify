package repository

type Repository[T any] interface {
	Create(entity *T) error
	Update(entity *T) error
	Delete(id string) error
	FindById(id string) (*T, error)
	FindAll(page int, limit int) (*PaginatedEntities[T], error)
}

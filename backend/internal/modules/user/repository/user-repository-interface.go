package repository

import "github.com/gianghp/statify/internal/database/models"

type IUserRepository interface {
	FindByID(id uint) (*models.User, error)
	FindAll() ([]*models.User, error)
	FindByEmail(email string) (*models.User, error)
	FindByUsername(username string) (*models.User, error)
	Create(user *models.User) error
	Update(user *models.User) error
	Delete(user *models.User) error
}

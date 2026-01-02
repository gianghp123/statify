package repository

import (
	"github.com/gianghp/statify/internal/database/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) FindByID(id uint) (*models.User, error) {
	return nil, nil
}

func (r *UserRepository) FindAll() ([]*models.User, error) {
	return nil, nil
}

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	return nil, nil
}

func (r *UserRepository) FindByUsername(username string) (*models.User, error) {
	return nil, nil
}

func (r *UserRepository) Create(user *models.User) error {
	return nil
}

func (r *UserRepository) Update(user *models.User) error {
	return nil
}

func (r *UserRepository) Delete(user *models.User) error {
	return nil
}

package service

import "github.com/gianghp/statify/internal/modules/user/repository"

type UserService struct {
	repo repository.IUserRepository
}

func NewUserService(repo repository.IUserRepository) *UserService {
	return &UserService{repo: repo}
}

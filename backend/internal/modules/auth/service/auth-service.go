package service

import (
	"github.com/gianghp/statify/internal/core"
	"github.com/gianghp/statify/internal/core/enums"
	"github.com/gianghp/statify/internal/database/models"
	"github.com/gianghp/statify/internal/modules/auth/dtos/request"
	"github.com/gianghp/statify/internal/modules/user/repository"
	"github.com/gianghp/statify/internal/utils"
	"github.com/gianghp/statify/internal/utils/bcrypt"
)

type AuthResult struct {
	User  *models.User
	Token string
}

type AuthService struct {
	userRepo repository.IUserRepository
	bcrypt   bcrypt.IBcryptUtils
	jwt      IJwtService
}

func NewAuthService(userRepo repository.IUserRepository, bcrypt bcrypt.IBcryptUtils, jwt IJwtService) *AuthService {
	return &AuthService{userRepo: userRepo, bcrypt: bcrypt, jwt: jwt}
}

func (s *AuthService) Register(request *request.RegisterRequest) (*models.User, error) {

	if request.Username == "" {
		return nil, core.BadRequestError("Username is required")
	}

	if err := utils.ValidateEmail(request.Email); err != nil {
		return nil, core.BadRequestError(err.Error())
	}

	if err := utils.ValidatePassword(request.Password); err != nil {
		return nil, core.BadRequestError(err.Error())
	}

	if _, err := s.userRepo.FindByEmail(request.Email); err != nil {
		return nil, core.ConflictError("Email already exists")
	}

	passwordHash, err := s.bcrypt.HashPassword(request.Password)
	if err != nil {
		return nil, core.InternalError()
	}

	user := &models.User{
		Email:        request.Email,
		PasswordHash: passwordHash,
		Username:     request.Username,
		Role:         string(enums.UserRoleUser),
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, core.InternalError()
	}

	return user, nil
}

func (s *AuthService) Login(request *request.LoginRequest) (*AuthResult, error) {
	user, err := s.userRepo.FindByEmail(request.Email)
	if err != nil {
		return nil, core.NotFoundError()
	}

	if err := s.bcrypt.CheckPassword(request.Password, user.PasswordHash); err != nil {
		return nil, core.UnauthorizedError()
	}

	token, err := s.jwt.Generate(user.ID, user.Role)
	if err != nil {
		return nil, core.InternalError()
	}

	return &AuthResult{User: user, Token: token}, nil
}

func (s *AuthService) Me(userID uint) (*models.User, error) {
	// TODO: implement
	return nil, nil
}

// func (s *AuthService) GitHubLogin() (*models.User, error) {
// 	// TODO: implement
// 	return nil, nil
// }

// func (s *AuthService) GitHubCallback() {
// 	// TODO: implement
// }

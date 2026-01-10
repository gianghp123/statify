package service

import (
	"context"
	"log"

	"github.com/gianghp/statify/internal/core"
	"github.com/gianghp/statify/internal/core/enums"
	"github.com/gianghp/statify/internal/database/models"
	"github.com/gianghp/statify/internal/modules/auth/dtos/request"
	userResponse "github.com/gianghp/statify/internal/modules/user/dtos/response"
	"github.com/gianghp/statify/internal/modules/user/repository"
	"github.com/gianghp/statify/internal/utils"
	"github.com/gianghp/statify/internal/utils/bcrypt"
)

type AuthResult struct {
	User  *userResponse.UserDto `json:"user"`
	Token string                `json:"token"`
}

type AuthService struct {
	userRepo repository.IUserRepository
	bcrypt   bcrypt.IBcryptUtils
	jwt      IJwtService
}

func NewAuthService(userRepo repository.IUserRepository, bcrypt bcrypt.IBcryptUtils, jwt IJwtService) *AuthService {
	return &AuthService{userRepo: userRepo, bcrypt: bcrypt, jwt: jwt}
}

func (s *AuthService) Register(ctx context.Context, request *request.RegisterRequest) (*userResponse.UserDto, error) {

	if request.Username == "" {
		return nil, core.BadRequestError("Username is required")
	}

	if err := utils.ValidateEmail(request.Email); err != nil {
		return nil, core.BadRequestError(err.Error())
	}

	if err := utils.ValidatePassword(request.Password); err != nil {
		return nil, core.BadRequestError(err.Error())
	}

	if existingUser, _ := s.userRepo.FindByEmail(ctx, request.Email); existingUser != nil {
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

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, core.InternalError()
	}

	userDto, err := utils.EntityToDto[userResponse.UserDto](user)
	if err != nil {
		return nil, core.InternalError()
	}

	return userDto, nil
}

func (s *AuthService) Login(ctx context.Context, request *request.LoginRequest) (*AuthResult, error) {
	user, err := s.userRepo.FindByEmail(ctx, request.Email)
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

	userDto, err := utils.EntityToDto[userResponse.UserDto](user)
	if err != nil {
		return nil, core.InternalError()
	}

	return &AuthResult{User: userDto, Token: token}, nil
}

func (s *AuthService) Me(ctx context.Context, userID uint) (*userResponse.UserDto, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, core.NotFoundError()
	}

	userDto, err := utils.EntityToDto[userResponse.UserDto](user)
	if err != nil {
		return nil, core.InternalError()
	}

	return userDto, nil
}

// func (s *AuthService) GitHubLogin() (*models.User, error) {
// 	// TODO: implement
// 	return nil, nil
// }

// func (s *AuthService) GitHubCallback() {
// 	// TODO: implement
// }

func CreateAdmin(userRepo repository.IUserRepository, bcrypt bcrypt.IBcryptUtils) error {
	if _, err := userRepo.FindByEmail(context.TODO(), "admin@admin.com"); err == nil {
		log.Println("Admin already exists")
		return nil
	}

	passwordHash, err := bcrypt.HashPassword("admin")
	if err != nil {
		return core.InternalError()
	}

	user := &models.User{
		Email:        "admin@admin.com",
		PasswordHash: passwordHash,
		Username:     "admin",
		Role:         string(enums.UserRoleAdmin),
	}

	if err := userRepo.Create(context.TODO(), user); err != nil {
		log.Println("Admin created failed")
		return core.InternalError()
	}

	log.Println("Admin created successfully")

	return nil
}

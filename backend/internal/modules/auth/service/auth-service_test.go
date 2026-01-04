package service

import (
	"context"
	"errors"
	"testing"

	"github.com/gianghp/statify/internal/core/enums"
	"github.com/gianghp/statify/internal/database/models"
	"github.com/gianghp/statify/internal/modules/auth/dtos/request"
	userResponse "github.com/gianghp/statify/internal/modules/user/dtos/response"
	"github.com/gianghp/statify/internal/modules/user/repository"
	"github.com/gianghp/statify/internal/utils/bcrypt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestAuthService_Register(t *testing.T) {
	tests := []struct {
		name         string
		request      *request.RegisterRequest
		setupMocks   func(repo *repository.UserRepositoryMock, bcrypt *bcrypt.BcryptUtilsMock)
		expectedFunc func(t *testing.T, user *userResponse.UserDto, err error)
	}{
		{
			name: "Register successfully",
			request: &request.RegisterRequest{
				Email:    "test@gmail.com",
				Password: "password",
				Username: "test",
			},
			setupMocks: func(repo *repository.UserRepositoryMock, bcrypt *bcrypt.BcryptUtilsMock) {
				repo.On("FindByEmail", mock.Anything, "test@gmail.com").Return((*models.User)(nil), nil)
				bcrypt.On("HashPassword", "password").Return("hashed-password", nil)
				repo.On("Create", mock.Anything, mock.MatchedBy(func(u *models.User) bool {
					return u.Email == "test@gmail.com" &&
						u.Username == "test" &&
						u.PasswordHash == "hashed-password" &&
						u.Role == string(enums.UserRoleUser)
				})).Return(nil)
			},
			expectedFunc: func(t *testing.T, user *userResponse.UserDto, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, "test@gmail.com", user.Email)
				assert.Equal(t, "test", user.Username)
				assert.Equal(t, string(enums.UserRoleUser), user.Role)
			},
		},
		{
			name: "Register with invalid email",
			request: &request.RegisterRequest{
				Email:    "test",
				Password: "password",
				Username: "test",
			},
			setupMocks: nil,
			expectedFunc: func(t *testing.T, user *userResponse.UserDto, err error) {
				assert.Error(t, err)
				assert.Nil(t, user)
			},
		},
		{
			name: "Register with invalid password",
			request: &request.RegisterRequest{
				Email:    "test@gmail.com",
				Password: "pass",
				Username: "test",
			},
			setupMocks: nil,
			expectedFunc: func(t *testing.T, user *userResponse.UserDto, err error) {
				assert.Error(t, err)
				assert.Nil(t, user)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(repository.UserRepositoryMock)
			bcrypt := new(bcrypt.BcryptUtilsMock)

			if tt.setupMocks != nil {
				tt.setupMocks(repo, bcrypt)
			}

			authService := NewAuthService(repo, bcrypt, nil)
			user, err := authService.Register(context.TODO(), tt.request)

			tt.expectedFunc(t, user, err)

			repo.AssertExpectations(t)
			bcrypt.AssertExpectations(t)
		})
	}
}

func TestAuthService_Login(t *testing.T) {
	tests := []struct {
		name         string
		request      *request.LoginRequest
		setupMocks   func(repo *repository.UserRepositoryMock, bcrypt *bcrypt.BcryptUtilsMock, jwt *JwtServiceMock)
		expectedFunc func(t *testing.T, authResult *AuthResult, err error)
	}{
		{
			name: "Login successfully",
			request: &request.LoginRequest{
				Email:    "test@gmail.com",
				Password: "password",
			},
			setupMocks: func(repo *repository.UserRepositoryMock, bcrypt *bcrypt.BcryptUtilsMock, jwt *JwtServiceMock) {
				repo.On("FindByEmail", mock.Anything, "test@gmail.com").Return(&models.User{
					Model: gorm.Model{
						ID: 1,
					},
					Email:        "test@gmail.com",
					PasswordHash: "hashed-password",
					Username:     "test",
					Role:         string(enums.UserRoleUser),
				}, nil)
				bcrypt.On("CheckPassword", "password", "hashed-password").Return(nil)
				jwt.On("Generate", uint(1), string(enums.UserRoleUser)).Return("token", nil)
			},
			expectedFunc: func(t *testing.T, authResult *AuthResult, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, authResult)
				assert.Equal(t, "test@gmail.com", authResult.User.Email)
				assert.Equal(t, "test", authResult.User.Username)
				assert.Equal(t, string(enums.UserRoleUser), authResult.User.Role)
				assert.Equal(t, "token", authResult.Token)
			},
		},
		{
			name: "Login with invalid email",
			request: &request.LoginRequest{
				Email:    "test",
				Password: "password",
			},
			setupMocks: func(repo *repository.UserRepositoryMock, bcrypt *bcrypt.BcryptUtilsMock, jwt *JwtServiceMock) {
				repo.On("FindByEmail", mock.Anything, "test").Return((*models.User)(nil), gorm.ErrRecordNotFound)
			},
			expectedFunc: func(t *testing.T, authResult *AuthResult, err error) {
				assert.Error(t, err)
				assert.Nil(t, authResult)
			},
		},
		{
			name: "Login with invalid password",
			request: &request.LoginRequest{
				Email:    "test@gmail.com",
				Password: "pass",
			},
			setupMocks: func(repo *repository.UserRepositoryMock, bcrypt *bcrypt.BcryptUtilsMock, jwt *JwtServiceMock) {
				repo.On("FindByEmail", mock.Anything, "test@gmail.com").Return(&models.User{
					Model: gorm.Model{
						ID: 1,
					},
					Email:        "test@gmail.com",
					PasswordHash: "hashed-password",
					Username:     "test",
					Role:         string(enums.UserRoleUser),
				}, nil)
				bcrypt.On("CheckPassword", "pass", "hashed-password").Return(errors.New("invalid password"))
			},
			expectedFunc: func(t *testing.T, authResult *AuthResult, err error) {
				assert.Error(t, err)
				assert.Nil(t, authResult)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(repository.UserRepositoryMock)
			bcrypt := new(bcrypt.BcryptUtilsMock)
			jwt := new(JwtServiceMock)

			if tt.setupMocks != nil {
				tt.setupMocks(repo, bcrypt, jwt)
			}

			authService := NewAuthService(repo, bcrypt, jwt)
			authResult, err := authService.Login(context.TODO(), tt.request)

			tt.expectedFunc(t, authResult, err)

			repo.AssertExpectations(t)
			bcrypt.AssertExpectations(t)
			jwt.AssertExpectations(t)
		})
	}
}

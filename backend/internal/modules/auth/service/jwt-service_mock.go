package service

import (
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/mock"
)

type JwtServiceMock struct {
	mock.Mock
}

func (j *JwtServiceMock) Generate(userID uint, role string) (string, error) {
	args := j.Called(userID, role)
	return args.String(0), args.Error(1)
}

func (j *JwtServiceMock) Verify(tokenString string) (*jwt.MapClaims, error) {
	args := j.Called(tokenString)
	return args.Get(0).(*jwt.MapClaims), args.Error(1)
}

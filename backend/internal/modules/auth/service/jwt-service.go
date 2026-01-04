package service

import (
	"time"

	"github.com/gianghp/statify/internal/core"
	"github.com/gianghp/statify/internal/utils"
	jwt "github.com/golang-jwt/jwt/v5"
)

type IJwtService interface {
	Generate(userID uint, role string) (string, error)
	Verify(tokenString string) (*jwt.MapClaims, error)
}

type JwtService struct {
	secretKey string
}

func NewJwtService() *JwtService {
	secretKey := utils.GetEnv("JWT_SECRET", "")
	return &JwtService{secretKey: secretKey}
}

func (j *JwtService) Generate(userID uint, role string) (string, error) {
	claims := jwt.MapClaims{
		"sub":  userID,
		"role": role,
		"exp":  time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secretKey))
}

func (j *JwtService) Verify(tokenString string) (*jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return []byte(j.secretKey), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, core.UnauthorizedError()
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return &claims, nil
	} else {
		return nil, core.UnauthorizedError()
	}
}

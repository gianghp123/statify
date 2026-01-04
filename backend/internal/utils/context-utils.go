package utils

import (
	"errors"

	"github.com/gin-gonic/gin"
)

func GetUserIDFromContext(ctx *gin.Context) (uint, error) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		return 0, errors.New("user_id not found in context")
	}

	uintUserID, ok := userID.(uint)
	if !ok {
		return 0, errors.New("user_id is not uint")
	}

	return uintUserID, nil
}

func GetRoleFromContext(ctx *gin.Context) (string, error) {
	role, exists := ctx.Get("role")
	if !exists {
		return "", errors.New("role not found in context")
	}

	return role.(string), nil
}

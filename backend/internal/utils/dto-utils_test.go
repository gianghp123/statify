package utils

import (
	"testing"
	"time"

	"github.com/gianghp/statify/internal/database/models"
	"github.com/gianghp/statify/internal/modules/user/dtos/response"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestEntityToDto(t *testing.T) {
	fixedTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	entity := &models.User{
		Model: gorm.Model{
			ID:        1,
			CreatedAt: fixedTime,
			UpdatedAt: fixedTime,
		},
		Username: "test",
		Email:    "test",
		Role:     "test",
	}

	expected := &response.UserDto{
		ID:        1,
		Username:  "test",
		Email:     "test",
		Role:      "test",
		CreatedAt: fixedTime,
		UpdatedAt: fixedTime,
	}

	dto, err := EntityToDto[*response.UserDto](&entity)
	assert.NoError(t, err)
	assert.Equal(t, expected, dto)
}

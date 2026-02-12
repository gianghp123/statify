package transaction

import (
	"context"

	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type TransactionManagerMock struct {
	mock.Mock
}

func (m *TransactionManagerMock) Transaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
	args := m.Called(ctx, fn)
	if rf, ok := args.Get(0).(func(context.Context, func(*gorm.DB) error) error); ok {
		return rf(ctx, fn)
	}
	return args.Error(0)
}

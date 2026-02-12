package transaction

import (
	"context"

	"gorm.io/gorm"
)

type TransactionManager struct {
	db *gorm.DB
}

func NewTransactionManager(db *gorm.DB) *TransactionManager {
	return &TransactionManager{db: db}
}

func (t *TransactionManager) Transaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
	return t.db.WithContext(ctx).Transaction(fn)
}

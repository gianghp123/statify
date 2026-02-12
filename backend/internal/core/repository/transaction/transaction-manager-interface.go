package transaction

import (
	"context"

	"gorm.io/gorm"
)

type ITransactionManager interface {
	Transaction(ctx context.Context, fn func(tx *gorm.DB) error) error
}

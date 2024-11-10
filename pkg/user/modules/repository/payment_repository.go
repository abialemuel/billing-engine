package repository

import (
	"context"

	"github.com/abialemuel/poly-kit/infrastructure/apm"
	"gorm.io/gorm"
)

func (r *PgDBRepository) CreatePayment(ctx context.Context, tx *gorm.DB, payment *Payment) error {
	ctx, span := apm.StartTransaction(ctx, "Repository::CreatePayment")
	defer apm.EndTransaction(span)

	db := r.db
	if tx != nil {
		db = tx
	}

	if err := db.WithContext(ctx).Create(payment).Error; err != nil {
		return err
	}
	return nil
}

package repository

import (
	"context"

	"github.com/abialemuel/poly-kit/infrastructure/apm"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (r *PgDBRepository) GetLoanByUserID(ctx context.Context, tx *gorm.DB, userID uint) (*Loan, error) {
	ctx, span := apm.StartTransaction(ctx, "Repository::GetLoanByUserID")
	defer apm.EndTransaction(span)

	db := r.db
	if tx != nil {
		db = tx
	}

	var loan Loan
	if err := db.WithContext(ctx).Clauses(clause.Locking{Strength: "UPDATE"}).Where("user_id = ?", userID).First(&loan).Error; err != nil {
		return nil, err
	}
	return &loan, nil
}

// UpdateLoan updates a existing loan record, using tx if provided
func (r *PgDBRepository) UpdateLoan(ctx context.Context, tx *gorm.DB, loan *Loan) error {
	ctx, span := apm.StartTransaction(ctx, "Repository::UpdateLoan")
	defer apm.EndTransaction(span)

	db := r.db
	if tx != nil {
		db = tx
	}

	if err := db.WithContext(ctx).Save(loan).Error; err != nil {
		return err
	}
	return nil
}

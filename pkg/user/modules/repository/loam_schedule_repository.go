package repository

import (
	"context"
	"time"

	"github.com/abialemuel/poly-kit/infrastructure/apm"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (r *PgDBRepository) GetOverdueScheduleByLoanID(ctx context.Context, tx *gorm.DB, loanID uint) ([]LoanSchedule, error) {
	ctx, span := apm.StartTransaction(ctx, "Repository::GetOverdueScheduleByLoanID")
	defer apm.EndTransaction(span)

	db := r.db
	if tx != nil {
		db = tx
	}

	var schedules []LoanSchedule
	if err := db.WithContext(ctx).Where("loan_id = ? AND due_date < ? AND is_paid = ?", loanID, time.Now(), false).Find(&schedules).Error; err != nil {
		return nil, err
	}

	return schedules, nil
}

func (r *PgDBRepository) CreateSchedule(ctx context.Context, tx *gorm.DB, schedule *LoanSchedule) error {
	ctx, span := apm.StartTransaction(ctx, "Repository::CreateSchedule")
	defer apm.EndTransaction(span)

	db := r.db
	if tx != nil {
		db = tx
	}

	if err := db.WithContext(ctx).Create(schedule).Error; err != nil {
		return err
	}
	return nil
}

// UpdateSchedule updates a existing schedule record use map[string]interface{} to update the record
func (r *PgDBRepository) UpdateSchedule(ctx context.Context, tx *gorm.DB, schedule *LoanSchedule) error {
	ctx, span := apm.StartTransaction(ctx, "Repository::UpdateSchedule")
	defer apm.EndTransaction(span)

	db := r.db
	if tx != nil {
		db = tx
	}

	if err := db.WithContext(ctx).Clauses(clause.Locking{Strength: "UPDATE"}).Save(schedule).Error; err != nil {
		return err
	}

	return nil
}

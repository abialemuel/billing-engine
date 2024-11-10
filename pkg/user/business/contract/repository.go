package contract

import (
	"context"

	"github.com/abialemuel/billing-engine/pkg/user/modules/repository"
	"gorm.io/gorm"
)

type Repository interface {
	// Users Repository
	GetUserByUsername(ctx context.Context, tx *gorm.DB, username string) (*repository.User, error)

	// Loans Repository
	GetLoanByUserID(ctx context.Context, tx *gorm.DB, userID uint) (*repository.Loan, error)
	UpdateLoan(ctx context.Context, tx *gorm.DB, loan *repository.Loan) error

	// Payments Repository
	CreatePayment(ctx context.Context, tx *gorm.DB, payment *repository.Payment) error

	// Loan Schedules Repository
	GetOverdueScheduleByLoanID(ctx context.Context, tx *gorm.DB, loanID uint) ([]repository.LoanSchedule, error)
	CreateSchedule(ctx context.Context, tx *gorm.DB, schedule *repository.LoanSchedule) error
	UpdateSchedule(ctx context.Context, tx *gorm.DB, schedule *repository.LoanSchedule) error
}

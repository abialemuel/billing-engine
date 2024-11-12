package business

import (
	"context"
	"fmt"

	mainCfg "github.com/abialemuel/billing-engine/config"
	"github.com/abialemuel/billing-engine/pkg/user/business/contract"
	"github.com/abialemuel/billing-engine/pkg/user/business/model"
	"github.com/abialemuel/billing-engine/pkg/user/modules/repository"
	"github.com/abialemuel/poly-kit/infrastructure/apm"
	"gorm.io/gorm"
)

// UserService Business Logic of user domain
type UserService struct {
	db   *gorm.DB
	repo contract.Repository
	cfg  mainCfg.Config
}

// NewUserService creates a new instance of UserService
func NewUserService(
	repo contract.Repository,
	cfg mainCfg.Config,
	db *gorm.DB,
) UserService {
	return UserService{repo: repo, cfg: cfg, db: db}
}

// GetOutstandingLoad
func (s *UserService) GetOutstandingLoan(ctx context.Context, userID string) (res model.OutstandingLoan, err error) {
	ctx, span := apm.StartTransaction(ctx, "UserService::GetOutstandingLoan")
	defer apm.EndTransaction(span)

	// Get user by ID
	user, err := s.repo.GetUserByUsername(ctx, s.db, userID)
	if err != nil {
		return res, err
	}

	// Get all loans by user ID
	loan, err := s.repo.GetLoanByUserID(ctx, s.db, user.ID)
	if err != nil {
		return res, err
	}

	// get all loan schedules by loan ID and check if the loan is overdue
	schedules, err := s.repo.GetOverdueScheduleByLoanID(ctx, s.db, loan.ID)
	if err != nil {
		return res, err
	}

	res.LoanID = loan.ID
	res.Username = user.Username
	res.OutstandingAmount = loan.Outstanding

	overdueSchedulesCount := uint(len(schedules))
	if overdueSchedulesCount >= uint(s.cfg.Get().Billing.DelinquentThreshold) {
		res.IsDelinquent = true
	} else {
		res.IsDelinquent = false
	}

	res.MissedPayment = overdueSchedulesCount
	res.UpcomingAmount = totalOverdueAmount(schedules)

	return res, nil
}

// MakePayment
func (s *UserService) MakePayment(ctx context.Context, userID string, amount float64) error {
	ctx, span := apm.StartTransaction(ctx, "UserService::MakePayment")
	defer apm.EndTransaction(span)

	// Start a transaction
	tx := s.db.Begin()
	if err := tx.Error; err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	// Ensure that transaction is either committed or rolled back
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r) // Re-panic after rollback to maintain original error
		} else if err := tx.Commit().Error; err != nil {
			tx.Rollback() // Rollback if commit fails
		}
	}()

	// Get user by ID
	user, err := s.repo.GetUserByUsername(ctx, tx, userID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Get the loan by user ID
	loan, err := s.repo.GetLoanByUserID(ctx, tx, user.ID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to get loan: %w", err)
	}

	// Get overdue loan schedules
	schedules, err := s.repo.GetOverdueScheduleByLoanID(ctx, tx, loan.ID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to get overdue schedules: %w", err)
	}

	// Validate overdue schedules and payment amount
	if len(schedules) == 0 {
		tx.Rollback()
		return fmt.Errorf("no overdue schedule found")
	}

	if totalOverdueAmount(schedules) != amount {
		tx.Rollback()
		return fmt.Errorf("amount should be equal to due amount")
	}

	// Update loan's outstanding amount
	loan.Outstanding -= amount
	if err := s.repo.UpdateLoan(ctx, tx, loan); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update loan: %w", err)
	}

	// Mark schedules as paid
	for _, schedule := range schedules {
		schedule.IsPaid = true
		if err := s.repo.UpdateSchedule(ctx, tx, &schedule); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to update loan schedule: %w", err)
		}
	}

	// create payment record
	payment := repository.Payment{
		LoanID:     loan.ID,
		AmountPaid: amount,
	}
	if err := s.repo.CreatePayment(ctx, tx, &payment); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create payment: %w", err)
	}

	return nil // Transaction will be committed by the deferred function
}

func totalOverdueAmount(schedules []repository.LoanSchedule) float64 {
	var total float64
	for _, schedule := range schedules {
		total += schedule.AmountDue
	}
	return total
}

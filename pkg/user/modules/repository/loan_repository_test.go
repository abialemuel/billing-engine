package repository_test

import (
	"context"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/abialemuel/billing-engine/pkg/user/modules/repository"
	"github.com/stretchr/testify/assert"
)

func (suite *RepositoryTestSuite) TestGetLoanByUserID_Success() {
	// Define test data
	userID := uint(1)
	expectedLoan := &repository.Loan{
		ID:     1,
		UserID: userID,
	}

	// Mock database query with SELECT
	suite.sqlMock.ExpectQuery(`SELECT \* FROM "loans" WHERE user_id = \$1 ORDER BY "loans"\."id" LIMIT \$2 FOR UPDATE`).
		WithArgs(userID, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id"}).AddRow(expectedLoan.ID, expectedLoan.UserID))

	// Call the function
	ctx := context.Background()
	loan, err := suite.repository.GetLoanByUserID(ctx, nil, userID)

	// Check the result
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedLoan, loan)
}

func (suite *RepositoryTestSuite) TestUpdateLoan_Success() {
	// Define test data
	loan := &repository.Loan{
		ID:           1,
		UserID:       1,
		Principal:    1000,
		InterestRate: 5.0,
		TotalWeeks:   52,
		WeeklyAmount: 20,
		Outstanding:  500,
		MissedAmount: 0,
		Delinquent:   nil,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Mock transaction start
	suite.sqlMock.ExpectBegin()

	// Mock database query with UPDATE
	suite.sqlMock.ExpectExec(`UPDATE "loans" SET "user_id"=\$1,"principal"=\$2,"interest_rate"=\$3,"total_weeks"=\$4,"weekly_amount"=\$5,"outstanding"=\$6,"missed_amount"=\$7,"delinquent"=\$8,"created_at"=\$9,"updated_at"=\$10 WHERE "id" = \$11`).
		WithArgs(
			loan.UserID,
			loan.Principal,
			loan.InterestRate,
			loan.TotalWeeks,
			loan.WeeklyAmount,
			loan.Outstanding,
			loan.MissedAmount,
			loan.Delinquent,
			sqlmock.AnyArg(), // Use AnyArg() for created_at
			sqlmock.AnyArg(), // Use AnyArg() for updated_at
			loan.ID,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Mock transaction commit
	suite.sqlMock.ExpectCommit()

	// Call the function
	ctx := context.Background()
	err := suite.repository.UpdateLoan(ctx, nil, loan)

	// Check the result
	assert.NoError(suite.T(), err)
}

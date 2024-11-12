package business_test

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/abialemuel/billing-engine/config"
	mockCfg "github.com/abialemuel/billing-engine/mocks/config"
	mockRepository "github.com/abialemuel/billing-engine/mocks/pkg/user/business/contract"
	"github.com/abialemuel/billing-engine/pkg/user/business"
	"github.com/abialemuel/billing-engine/pkg/user/modules/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Suite struct {
	suite.Suite
	controller *gomock.Controller
	mockDB     *gorm.DB
	sqlMock    sqlmock.Sqlmock
}

func (c *Suite) SetupTest() {
	c.controller = gomock.NewController(c.T())

	// Initialize sqlmock
	db, sqlMock, err := sqlmock.New()
	if err != nil {
		c.T().Fatalf("error initializing sqlmock: %v", err)
	}
	c.sqlMock = sqlMock

	// Initialize gorm DB with sqlmock
	c.mockDB, err = gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		c.T().Fatalf("error initializing gorm with sqlmock: %v", err)
	}
}

func (c *Suite) TearDownTest() {
	c.controller.Finish()
}

func (c *Suite) TestGetOutstandingLoan() {
	mockRepo := mockRepository.NewMockRepository(c.controller)
	mockConfig := mockCfg.NewMockConfig(c.controller)
	mockConfig.EXPECT().Get().Return(&config.MainConfig{}).AnyTimes()

	service := business.NewUserService(mockRepo, mockConfig, c.mockDB)

	// Define test data
	userID := 1
	username := "test"

	expectedUser := &repository.User{
		ID:       uint(userID),
		Username: username,
	}
	expectedLoan := &repository.Loan{
		ID:          1,
		UserID:      uint(userID),
		Outstanding: 1000,
	}

	expectedSchedules := []repository.LoanSchedule{
		{
			LoanID:    uint(userID),
			AmountDue: 100,
		},
	}

	c.Run("GetOutstandingLoan success", func() {
		// Set up expectations
		mockRepo.EXPECT().GetUserByUsername(gomock.Any(), c.mockDB, username).Return(expectedUser, nil)
		mockRepo.EXPECT().GetLoanByUserID(gomock.Any(), c.mockDB, uint(userID)).Return(expectedLoan, nil)
		mockRepo.EXPECT().GetOverdueScheduleByLoanID(gomock.Any(), c.mockDB, uint(userID)).Return(expectedSchedules, nil)

		// Call the method
		ctx := context.Background()
		loan, err := service.GetOutstandingLoan(ctx, username)

		// Assert results
		assert.NoError(c.T(), err)
		assert.Equal(c.T(), expectedLoan.Outstanding, loan.OutstandingAmount)
	})

	c.Run("GetOutstandingLoan failed to get user", func() {
		// Set up expectations
		mockRepo.EXPECT().GetUserByUsername(gomock.Any(), c.mockDB, username).Return(nil, assert.AnError)

		// Call the method
		ctx := context.Background()
		_, err := service.GetOutstandingLoan(ctx, username)

		// Assert results
		assert.Error(c.T(), err)
	})

	c.Run("GetOutstandingLoan failed to get loan", func() {
		// Set up expectations
		mockRepo.EXPECT().GetUserByUsername(gomock.Any(), c.mockDB, username).Return(expectedUser, nil)
		mockRepo.EXPECT().GetLoanByUserID(gomock.Any(), c.mockDB, uint(userID)).Return(nil, assert.AnError)

		// Call the method
		ctx := context.Background()
		_, err := service.GetOutstandingLoan(ctx, username)

		// Assert results
		assert.Error(c.T(), err)
	})

	c.Run("GetOutstandingLoan failed to get schedules", func() {
		// Set up expectations
		mockRepo.EXPECT().GetUserByUsername(gomock.Any(), c.mockDB, username).Return(expectedUser, nil)
		mockRepo.EXPECT().GetLoanByUserID(gomock.Any(), c.mockDB, uint(userID)).Return(expectedLoan, nil)
		mockRepo.EXPECT().GetOverdueScheduleByLoanID(gomock.Any(), c.mockDB, uint(userID)).Return(nil, assert.AnError)

		// Call the method
		ctx := context.Background()
		_, err := service.GetOutstandingLoan(ctx, username)

		// Assert results
		assert.Error(c.T(), err)
	})

	c.Run("GetOutstandingLoan delinquent", func() {
		// Set up expectations
		mockRepo.EXPECT().GetUserByUsername(gomock.Any(), c.mockDB, username).Return(expectedUser, nil)
		mockRepo.EXPECT().GetLoanByUserID(gomock.Any(), c.mockDB, uint(userID)).Return(expectedLoan, nil)
		mockRepo.EXPECT().GetOverdueScheduleByLoanID(gomock.Any(), c.mockDB, uint(userID)).Return(expectedSchedules, nil)

		// Call the method
		ctx := context.Background()
		loan, err := service.GetOutstandingLoan(ctx, username)

		// Assert results
		assert.NoError(c.T(), err)
		assert.True(c.T(), loan.IsDelinquent)
	})

	c.Run("GetOutstandingLoan not delinquent", func() {
		// mockCfg to return Billing.DelinquentThreshold = 2
		mockConfig.Get().Billing.DelinquentThreshold = 2

		// Set up expectations
		mockRepo.EXPECT().GetUserByUsername(gomock.Any(), c.mockDB, username).Return(expectedUser, nil)
		mockRepo.EXPECT().GetLoanByUserID(gomock.Any(), c.mockDB, uint(userID)).Return(expectedLoan, nil)
		mockRepo.EXPECT().GetOverdueScheduleByLoanID(gomock.Any(), c.mockDB, uint(userID)).Return(expectedSchedules, nil)

		// Call the method
		ctx := context.Background()
		loan, err := service.GetOutstandingLoan(ctx, username)

		// Assert results
		assert.NoError(c.T(), err)
		assert.False(c.T(), loan.IsDelinquent)
	})
}

func (c *Suite) TestMakePayment() {
	mockRepo := mockRepository.NewMockRepository(c.controller)
	mockConfig := mockCfg.NewMockConfig(c.controller)
	mockConfig.EXPECT().Get().Return(&config.MainConfig{}).AnyTimes()

	service := business.NewUserService(mockRepo, mockConfig, c.mockDB)

	// Define test data
	userID := 1
	username := "test"
	amount := 100

	expectedUser := &repository.User{
		ID:       uint(userID),
		Username: username,
	}
	expectedLoan := &repository.Loan{
		ID:          1,
		UserID:      uint(userID),
		Outstanding: 1000,
	}

	expectedSchedules := []repository.LoanSchedule{
		{
			LoanID:    uint(userID),
			AmountDue: 100,
		},
	}

	// failed to start transaction
	c.Run("MakePayment failed to start transaction", func() {
		// Set up expectations
		c.sqlMock.ExpectBegin().WillReturnError(assert.AnError)

		// Call the method
		ctx := context.Background()
		err := service.MakePayment(ctx, username, float64(amount))

		// Assert results
		assert.Error(c.T(), err)
	})

	c.Run("MakePayment success", func() {
		// Set up expectations
		c.sqlMock.ExpectBegin()
		c.sqlMock.ExpectCommit()

		mockRepo.EXPECT().GetUserByUsername(gomock.Any(), gomock.Any(), username).Return(expectedUser, nil)
		mockRepo.EXPECT().GetLoanByUserID(gomock.Any(), gomock.Any(), uint(userID)).Return(expectedLoan, nil)
		mockRepo.EXPECT().GetOverdueScheduleByLoanID(gomock.Any(), gomock.Any(), uint(userID)).Return(expectedSchedules, nil)
		mockRepo.EXPECT().UpdateLoan(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
		mockRepo.EXPECT().UpdateSchedule(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
		mockRepo.EXPECT().CreatePayment(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

		// Call the method
		ctx := context.Background()
		err := service.MakePayment(ctx, username, float64(amount))

		// Assert results
		assert.NoError(c.T(), err)
	})

	c.Run("MakePayment failed to get user", func() {
		// Set up expectations
		c.sqlMock.ExpectBegin()
		c.sqlMock.ExpectRollback()

		mockRepo.EXPECT().GetUserByUsername(gomock.Any(), gomock.Any(), username).Return(nil, assert.AnError)

		// Call the method
		ctx := context.Background()
		err := service.MakePayment(ctx, username, float64(amount))

		// Assert results
		assert.Error(c.T(), err)
	})

	c.Run("MakePayment failed to get loan", func() {
		// Set up expectations
		c.sqlMock.ExpectBegin()
		c.sqlMock.ExpectRollback()

		mockRepo.EXPECT().GetUserByUsername(gomock.Any(), gomock.Any(), username).Return(expectedUser, nil)
		mockRepo.EXPECT().GetLoanByUserID(gomock.Any(), gomock.Any(), uint(userID)).Return(nil, assert.AnError)

		// Call the method
		ctx := context.Background()
		err := service.MakePayment(ctx, username, float64(amount))

		// Assert results
		assert.Error(c.T(), err)
	})

	c.Run("MakePayment failed to get schedules", func() {
		// Set up expectations
		c.sqlMock.ExpectBegin()
		c.sqlMock.ExpectRollback()

		mockRepo.EXPECT().GetUserByUsername(gomock.Any(), gomock.Any(), username).Return(expectedUser, nil)
		mockRepo.EXPECT().GetLoanByUserID(gomock.Any(), gomock.Any(), uint(userID)).Return(expectedLoan, nil)
		mockRepo.EXPECT().GetOverdueScheduleByLoanID(gomock.Any(), gomock.Any(), uint(userID)).Return(nil, assert.AnError)

		// Call the method
		ctx := context.Background()
		err := service.MakePayment(ctx, username, float64(amount))

		// Assert results
		assert.Error(c.T(), err)
	})

	c.Run("MakePayment no overdue schedules", func() {
		// Set up expectations
		c.sqlMock.ExpectBegin()
		c.sqlMock.ExpectRollback()

		mockRepo.EXPECT().GetUserByUsername(gomock.Any(), gomock.Any(), username).Return(expectedUser, nil)
		mockRepo.EXPECT().GetLoanByUserID(gomock.Any(), gomock.Any(), uint(userID)).Return(expectedLoan, nil)
		mockRepo.EXPECT().GetOverdueScheduleByLoanID(gomock.Any(), gomock.Any(), uint(userID)).Return([]repository.LoanSchedule{}, nil)

		// Call the method
		ctx := context.Background()
		err := service.MakePayment(ctx, username, float64(amount))

		// Assert results
		assert.Error(c.T(), err)
	})

	c.Run("MakePayment insufficient amount", func() {
		// Set up expectations
		c.sqlMock.ExpectBegin()
		c.sqlMock.ExpectRollback()

		mockRepo.EXPECT().GetUserByUsername(gomock.Any(), gomock.Any(), username).Return(expectedUser, nil)
		mockRepo.EXPECT().GetLoanByUserID(gomock.Any(), gomock.Any(), uint(userID)).Return(expectedLoan, nil)
		mockRepo.EXPECT().GetOverdueScheduleByLoanID(gomock.Any(), gomock.Any(), uint(userID)).Return(expectedSchedules, nil)

		// Call the method
		ctx := context.Background()
		err := service.MakePayment(ctx, username, 50)

		// Assert results
		assert.Error(c.T(), err)
	})

	c.Run("MakePayment failed to update loan", func() {
		// Set up expectations
		c.sqlMock.ExpectBegin()
		c.sqlMock.ExpectRollback()

		mockRepo.EXPECT().GetUserByUsername(gomock.Any(), gomock.Any(), username).Return(expectedUser, nil)
		mockRepo.EXPECT().GetLoanByUserID(gomock.Any(), gomock.Any(), uint(userID)).Return(expectedLoan, nil)
		mockRepo.EXPECT().GetOverdueScheduleByLoanID(gomock.Any(), gomock.Any(), uint(userID)).Return(expectedSchedules, nil)
		mockRepo.EXPECT().UpdateLoan(gomock.Any(), gomock.Any(), gomock.Any()).Return(assert.AnError)

		// Call the method
		ctx := context.Background()
		err := service.MakePayment(ctx, username, float64(amount))

		// Assert results
		assert.Error(c.T(), err)
	})

	c.Run("MakePayment failed to update schedule", func() {
		// Set up expectations
		c.sqlMock.ExpectBegin()
		c.sqlMock.ExpectRollback()

		mockRepo.EXPECT().GetUserByUsername(gomock.Any(), gomock.Any(), username).Return(expectedUser, nil)
		mockRepo.EXPECT().GetLoanByUserID(gomock.Any(), gomock.Any(), uint(userID)).Return(expectedLoan, nil)
		mockRepo.EXPECT().GetOverdueScheduleByLoanID(gomock.Any(), gomock.Any(), uint(userID)).Return(expectedSchedules, nil)
		mockRepo.EXPECT().UpdateLoan(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
		mockRepo.EXPECT().UpdateSchedule(gomock.Any(), gomock.Any(), gomock.Any()).Return(assert.AnError)

		// Call the method
		ctx := context.Background()
		err := service.MakePayment(ctx, username, float64(amount))

		// Assert results
		assert.Error(c.T(), err)
	})

	c.Run("MakePayment failed to create payment", func() {
		// Set up expectations
		c.sqlMock.ExpectBegin()
		c.sqlMock.ExpectRollback()

		mockRepo.EXPECT().GetUserByUsername(gomock.Any(), gomock.Any(), username).Return(expectedUser, nil)
		mockRepo.EXPECT().GetLoanByUserID(gomock.Any(), gomock.Any(), uint(userID)).Return(expectedLoan, nil)
		mockRepo.EXPECT().GetOverdueScheduleByLoanID(gomock.Any(), gomock.Any(), uint(userID)).Return(expectedSchedules, nil)
		mockRepo.EXPECT().UpdateLoan(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
		mockRepo.EXPECT().UpdateSchedule(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
		mockRepo.EXPECT().CreatePayment(gomock.Any(), gomock.Any(), gomock.Any()).Return(assert.AnError)

		// Call the method
		ctx := context.Background()
		err := service.MakePayment(ctx, username, float64(amount))

		// Assert results
		assert.Error(c.T(), err)
	})

}

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}

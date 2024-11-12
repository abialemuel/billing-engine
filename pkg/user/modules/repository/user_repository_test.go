package repository_test

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/abialemuel/billing-engine/pkg/user/modules/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type RepositoryTestSuite struct {
	suite.Suite
	mockCtrl   *gomock.Controller
	mockDB     *gorm.DB
	sqlMock    sqlmock.Sqlmock
	repository *repository.PgDBRepository
}

func (suite *RepositoryTestSuite) SetupTest() {
	suite.mockCtrl = gomock.NewController(suite.T())

	// Initialize sqlmock
	db, sqlMock, err := sqlmock.New()
	if err != nil {
		suite.T().Fatalf("error initializing sqlmock: %v", err)
	}
	suite.sqlMock = sqlMock

	// Initialize gorm DB with sqlmock
	suite.mockDB, err = gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		suite.T().Fatalf("error initializing gorm with sqlmock: %v", err)
	}

	// Initialize repository
	suite.repository = repository.NewPgDBRepository(suite.mockDB)
}

func (suite *RepositoryTestSuite) TearDownTest() {
	suite.mockCtrl.Finish()
}

func (suite *RepositoryTestSuite) TestGetUserByUsername_Success() {
	// Define test data
	username := "testuser"
	expectedUser := repository.User{
		ID:       1,
		Username: username,
	}

	// Mock database query with ORDER BY and LIMIT
	rows := sqlmock.NewRows([]string{"id", "username"}).AddRow(1, username)
	suite.sqlMock.ExpectQuery(`SELECT \* FROM "users" WHERE username = \$1 ORDER BY "users"."id" LIMIT \$2`).
		WithArgs(username, 1).
		WillReturnRows(rows)

	// Call the method
	ctx := context.Background()
	user, err := suite.repository.GetUserByUsername(ctx, nil, username)

	// Assert results
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), &expectedUser, user)
}

func TestRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}

package repository_test

// func (suite *RepositoryTestSuite) TestCreatePayment_Success() {
// 	// Define test data
// 	payment := &repository.Payment{
// 		ID:         1,
// 		LoanID:     1,
// 		AmountPaid: 1000,
// 		CreatedAt:  time.Now(),
// 	}

// 	// Mock transaction start
// 	suite.sqlMock.ExpectBegin()

// 	// Mock database query with INSERT
// 	suite.sqlMock.ExpectExec(`INSERT INTO "payments" \("loan_id","amount_paid","created_at","id"\) VALUES \(\$1,\$2,\$3,\$4\) RETURNING "id"`).
// 		WithArgs(payment.LoanID, payment.AmountPaid, sqlmock.AnyArg(), payment.ID).
// 		WillReturnResult(sqlmock.NewResult(1, 1))

// 	// Mock transaction commit
// 	suite.sqlMock.ExpectCommit()

// 	// Call the function
// 	ctx := context.Background()
// 	err := suite.repository.CreatePayment(ctx, nil, payment)

// 	// Check the result
// 	assert.NoError(suite.T(), err)
// }

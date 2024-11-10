package response

import "github.com/abialemuel/billing-engine/pkg/user/business/model"

type OutstandingLoanInfoResponse struct {
	Code    int                 `json:"code"`
	Message string              `json:"message"`
	Payload OutstandingLoanInfo `json:"payload"`
}

type OutstandingLoanInfo struct {
	LoanID         uint    `json:"loan_id"`
	Username       string  `json:"username"`
	Outstanding    float64 `json:"outstanding"`
	IsDelinquent   bool    `json:"is_delinquent"`
	UpcomingAmount float64 `json:"upcoming_amount"`
	MissedPayment  uint    `json:"missed_payment"`
}

func NewOutstandingLoanInfoResponse(info model.OutstandingLoan) OutstandingLoanInfoResponse {
	return OutstandingLoanInfoResponse{
		Code:    200,
		Message: "Success",
		Payload: OutstandingLoanInfo{
			LoanID:         info.LoanID,
			Username:       info.Username,
			Outstanding:    info.OutstandingAmount,
			IsDelinquent:   info.IsDelinquent,
			UpcomingAmount: info.UpcomingAmount,
			MissedPayment:  info.MissedPayment,
		},
	}

}

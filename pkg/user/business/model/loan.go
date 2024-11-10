package model

type Loan struct {
	ID          int     `json:"id"`
	UserID      int     `json:"user_id"`
	Principal   int     `json:"principal"`
	Interest    float64 `json:"interest"`
	TotalWeeks  int     `json:"total_weeks"`
	Weekly      int     `json:"weekly"`
	Outstanding int     `json:"outstanding"`
	Missed      int     `json:"missed"`
	Delinquent  bool    `json:"delinquent"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

type OutstandingLoan struct {
	LoanID            uint    `json:"loan_id"`
	Username          string  `json:"username"`
	OutstandingAmount float64 `json:"outstanding_amount"`
	IsDelinquent      bool    `json:"is_delinquent"`
	UpcomingAmount    float64 `json:"upcoming_amount"`
	MissedPayment     uint    `json:"missed_payment"`
}

type LoanSchedule struct {
	ID         int    `json:"id"`
	LoanID     int    `json:"loan_id"`
	WeekNumber int    `json:"week_number"`
	AmountDue  int    `json:"amount_due"`
	DueDate    string `json:"due_date"`
	IsPaid     bool   `json:"is_paid"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

type Payment struct {
	ID         int    `json:"id"`
	LoanID     int    `json:"loan_id"`
	AmountPaid int    `json:"amount_paid"`
	CreatedAt  string `json:"created_at"`
}

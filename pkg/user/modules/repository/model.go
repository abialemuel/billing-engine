package repository

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey"`
	Name      string    `gorm:"type:varchar(100);not null"`
	Username  string    `gorm:"type:varchar(100);unique"`
	Email     string    `gorm:"type:varchar(100);unique"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`

	Loans []Loan `gorm:"foreignKey:UserID"`
}

type Loan struct {
	ID           uint      `gorm:"primaryKey"`
	UserID       uint      `gorm:"not null"`
	Principal    int       `gorm:"not null"`
	InterestRate float64   `gorm:"type:decimal(5,2);not null"`
	TotalWeeks   int       `gorm:"not null"`
	WeeklyAmount int       `gorm:"not null"`
	Outstanding  float64   `gorm:"not null"`
	MissedAmount int       `gorm:"default:0"`
	Delinquent   *bool     `gorm:"default:false"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime"`

	User          User           `gorm:"foreignKey:UserID"`
	Payments      []Payment      `gorm:"foreignKey:LoanID"`
	LoanSchedules []LoanSchedule `gorm:"foreignKey:LoanID"`
}

type Payment struct {
	ID         uint      `gorm:"primaryKey"`
	LoanID     uint      `gorm:"not null"`
	AmountPaid float64   `gorm:"not null"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`

	Loan Loan `gorm:"foreignKey:LoanID"`
}

type LoanSchedule struct {
	ID         uint      `gorm:"primaryKey"`
	LoanID     uint      `gorm:"not null"`
	WeekNumber int       `gorm:"not null"`
	AmountDue  float64   `gorm:"not null"`
	DueDate    time.Time `gorm:"not null"`
	IsPaid     bool      `gorm:"default:false"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime"`

	Loan Loan `gorm:"foreignKey:LoanID"`
}

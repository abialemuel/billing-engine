package repository

import (
	"gorm.io/gorm"
)

// PgDBRepository The implementation of user.Repository object
type PgDBRepository struct {
	db *gorm.DB
}

// NewPgDBRepository Generate pg user repository
func NewPgDBRepository(db *gorm.DB) *PgDBRepository {
	repo := PgDBRepository{db: db}

	return &repo
}

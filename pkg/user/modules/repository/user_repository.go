package repository

import (
	"context"

	"github.com/abialemuel/poly-kit/infrastructure/apm"
	"gorm.io/gorm"
)

func (r *PgDBRepository) GetUserByUsername(ctx context.Context, tx *gorm.DB, username string) (*User, error) {
	ctx, span := apm.StartTransaction(ctx, "Repository::GetUserByUsername")
	defer apm.EndTransaction(span)

	db := r.db
	if tx != nil {
		db = tx
	}

	var user User
	err := db.WithContext(ctx).Where("username = ?", username).First(&user).Error
	return &user, err
}

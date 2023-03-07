package main

import (
	"context"

	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) PostUser(ctx context.Context, user User) error {
	return r.db.WithContext(ctx).Create(&user).Error
}

func (r *repository) GetUserWithEmail(ctx context.Context, email string) (user User, err error) {
	err = r.db.WithContext(ctx).Select("id, first_name, last_name, password").
		Find(&user, "email = ? AND active_status = 1", email).Error

	return
}

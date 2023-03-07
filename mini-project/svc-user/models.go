package main

import (
	"database/sql"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type InsertUser struct {
	Email           string `json:"email" validate:"email,max=100"`
	FirstName       string `json:"first_name" validate:"required,max=100"`
	LastName        string `json:"last_name" validate:"required,max=100"`
	Password        string `json:"password" validate:"required,max=128"`
	ConfirmPassword string `json:"confirm_password" validate:"eqfield=Password"`
}

func (i *InsertUser) RegisterUser() User {
	hash, _ := bcrypt.GenerateFromPassword([]byte(i.Password), bcrypt.MinCost)

	return User{
		Email:        i.Email,
		FirstName:    i.FirstName,
		LastName:     i.LastName,
		Password:     string(hash),
		ActiveStatus: 1,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}

type UserLogin struct {
	Email    string `json:"email" validate:"email"`
	Password string `json:"password" validate:"required"`
}

type User struct {
	ID           uint64       `gorm:"type:int unsigned;primaryKey;autoIncrement;not null"`
	Email        string       `gorm:"type:varchar(100);uniqueIndex:users_unique_index,priority:1;not null"`
	FirstName    string       `gorm:"type:varchar(100);not null"`
	LastName     string       `gorm:"type:varchar(100);not null"`
	Password     string       `gorm:"type:varchar(128);not null"`
	ActiveStatus uint8        `gorm:"type:tinyint unsigned;uniqueIndex:users_unique_index,priority:2;not null;default:1"`
	CreatedAt    time.Time    `gorm:"type:datetime;not null"`
	UpdatedAt    time.Time    `gorm:"type:datetime;not null"`
	DeletedAt    sql.NullTime `gorm:"type:datetime"`
}

type Response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

type PayloadNotify struct {
	Email     string `json:"email" validate:"required,email"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
}

package user

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name         string `gorm:"column:name"`
	Username     string `gorm:"username"`
	Password     string `gorm:"column:password"`
	PasswordSalt string `gorm:"column:password_salt"`
	Email        string `gorm:"column:email"`
	IsVerified   bool   `gorm:"column:is_verified"`
}

func (User) TableName() string {
	return "users"
}

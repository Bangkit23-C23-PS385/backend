package verification

import (
	"gorm.io/gorm"
)

type Verification struct {
	gorm.Model
	Email       string `gorm:"column:email"`
	Token       string `gorm:"column:token"`
	AttemptLeft int    `gorm:"column:attempt_left"`
}

package profile

import (
	"time"

	"ta/backend/src/constant"

	"gorm.io/gorm"
)

type Profile struct {
	gorm.Model
	UserId      string              `gorm:"column:userId"`
	Name        string              `gorm:"column:name"`
	Gender      constant.GenderType `gorm:"column:gender"`
	DateOfBirth time.Time           `gorm:"column:dateOfBirth"`
	Height      int                 `gorm:"column:height"`
	Weight      int                 `gorm:"column:weight"`
}

func (Profile) TableName() string {
	return "profile"
}

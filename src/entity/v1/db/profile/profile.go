package profile

import (
	"time"

	"backend/src/constant"
)

type Profile struct {
	// gorm.Model
	UserId      string              `gorm:"column:userid;primaryKey"`
	Name        string              `gorm:"column:name"`
	Gender      constant.GenderType `gorm:"column:gender"`
	DateOfBirth time.Time           `gorm:"column:dateofbirth"`
	Height      int                 `gorm:"column:height"`
	Weight      int                 `gorm:"column:weight"`
	CreatedAt   time.Time           `gorm:"column:created_at"`
	UpdatedAt   time.Time           `gorm:"column:updated_at"`
}

func (Profile) TableName() string {
	return "profiles"
}

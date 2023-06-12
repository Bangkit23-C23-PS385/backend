package profile

import (
	db "ta/backend/src/database"
	dbProfile "ta/backend/src/entity/v1/db/profile"

	"gorm.io/gorm"
)

type Repository struct {
	master *gorm.DB
}

func NewRepository(db db.DB) *Repository {
	return &Repository{
		master: db.Master,
	}
}

type Repositorier interface {
	CreateProfile(req dbProfile.Profile) (err error)
	GetProfile(id string) (entities dbProfile.Profile, err error)
	UpdateProfile(req dbProfile.Profile) (err error)
	DeleteProfile(id string) (err error)
}

func (repo Repository) CreateProfile(req dbProfile.Profile) (err error) {
	profileEntity := dbProfile.Profile{
		UserId:      req.UserId,
		Name:        req.Name,
		Gender:      req.Gender,
		DateOfBirth: req.DateOfBirth,
		Height:      req.Height,
		Weight:      req.Weight,
	}
	query := repo.master.Model(&dbProfile.Profile{}).Begin().Create(&profileEntity)
	err = query.Error
	if err != nil {
		query.Rollback()
		return
	}
	err = query.Commit().Error

	return
}

func (repo Repository) GetProfile(id string) (entities dbProfile.Profile, err error) {
	err = repo.master.
		Where("userId in (?)", id).
		Find(&entities).Error

	return
}

func (repo Repository) UpdateProfile(req dbProfile.Profile) (err error) {
	query := repo.master.Model(&dbProfile.Profile{}).Begin().
		Where("userId = ?", req.UserId).
		Updates(&req)
	err = query.Error
	if err != nil {
		query.Rollback()
		return
	}
	err = query.Commit().Error

	return
}

func (repo Repository) DeleteProfile(id string) (err error) {
	query := repo.master.Model(&dbProfile.Profile{}).Begin().Delete(&dbProfile.Profile{}, id)
	err = query.Error
	if err != nil {
		query.Rollback()
		return
	}
	err = query.Commit().Error

	return
}

package profile

import (
	db "backend/src/database"
	dbProfile "backend/src/entity/v1/db/profile"
	"log"

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
	GetProfile(userid string) (entities dbProfile.Profile, err error)
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

func (repo Repository) GetProfile(userid string) (entities dbProfile.Profile, err error) {
	err = repo.master.
		Where("userid in (?)", userid).
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

	// instance := repo.master.Model(dbProfile.Profile{})
	instance := repo.master.Model(dbProfile.Profile{UserId: id})
	result := instance.Delete(dbProfile.Profile{UserId: id})
	log.Println(result)

	if result.Error != nil {
		log.Println("failed to delete profile with id:" + id + " and error: " + result.Error.Error())
		return result.Error
	}

	return
}

package auth

import (
	"ta/backend/src/constant"
	db "ta/backend/src/database"
	dbUser "ta/backend/src/entity/v1/db/user"
	httpAuth "ta/backend/src/entity/v1/http/auth"

	"gorm.io/gorm"
)

type Repository struct {
	master *gorm.DB
	slave  *gorm.DB
}

func NewRepository(db db.DB) *Repository {
	return &Repository{
		master: db.Master,
		slave:  db.Slave,
	}
}

type Repositorier interface {
	GetUsers() (entities []dbUser.User, err error)
	GetUserByIDs(ids []int) (entities []dbUser.User, err error)
	GetUserByEmail(email string) (user dbUser.User, err error)
	Insert(req httpAuth.RegisterRequest) (err error)
	Update(userEntity dbUser.User) (err error)
	Delete(id int) (err error)
}

func (repo Repository) GetUsers() (entities []dbUser.User, err error) {
	err = repo.slave.
		Order("updated_at desc").
		Find(&entities).Error

	return
}

func (repo Repository) GetUserByIDs(ids []int) (entities []dbUser.User, err error) {
	err = repo.slave.
		Where("id in (?)", ids).
		Find(&entities).Error

	return
}

func (repo Repository) GetUserByEmail(email string) (user dbUser.User, err error) {
	err = repo.slave.
		Where("email = ?", email).
		Take(&user).Error

	return
}

func (repo Repository) Insert(req httpAuth.RegisterRequest) (err error) {
	userEntity := dbUser.User{
		Name:       req.Name,
		Email:      req.Email,
		IsVerified: false,
		Role:       constant.Roles(req.Role),
	}
	query := repo.master.Model(&dbUser.User{}).Begin().Create(&userEntity)
	err = query.Error
	if err != nil {
		query.Rollback()
		return
	}
	err = query.Commit().Error

	return
}

func (repo Repository) Update(userEntity dbUser.User) (err error) {
	query := repo.master.Model(&dbUser.User{}).Begin().
		Where("email = ?", userEntity.Email).
		Updates(&userEntity)
	err = query.Error
	if err != nil {
		query.Rollback()
		return
	}
	err = query.Commit().Error

	return
}

func (repo Repository) Delete(id int) (err error) {
	query := repo.master.Model(&dbUser.User{}).Begin().Delete(&dbUser.User{}, id)
	err = query.Error
	if err != nil {
		query.Rollback()
		return
	}
	err = query.Commit().Error

	return
}

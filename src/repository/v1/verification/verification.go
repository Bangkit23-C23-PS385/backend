package verification

import (
	db "ta/backend/src/database"
	dbVerification "ta/backend/src/entity/v1/db/verification"

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
	GetByEmail(email string) (verifEntity dbVerification.Verification, err error)
	Insert(email, token string) (err error)
	UpdateToken(email, token string, attemptLeft int) (err error)
	DeleteByEmail(email string) (err error)
}

func (repo Repository) GetByEmail(email string) (verifEntity dbVerification.Verification, err error) {
	err = repo.slave.
		Where("email = ?", email).
		Take(&verifEntity).Error

	return
}

func (repo Repository) Insert(email, token string) (err error) {
	verificationEntity := dbVerification.Verification{
		Email:       email,
		Token:       token,
		AttemptLeft: 3,
	}

	query := repo.master.Model(&dbVerification.Verification{}).Begin().Create(&verificationEntity)
	err = query.Error
	if err != nil {
		query.Rollback()
		return
	}
	err = query.Commit().Error

	return
}

func (repo Repository) UpdateToken(email, token string, attemptLeft int) (err error) {
	query := repo.master.Model(&dbVerification.Verification{}).Begin().
		Where("email = ?", email).
		Update("token", token).
		Update("attempt_left", attemptLeft)
	err = query.Error
	if err != nil {
		query.Rollback()
		return
	}
	err = query.Commit().Error

	return
}

func (repo Repository) DeleteByEmail(email string) (err error) {
	query := repo.master.Model(&dbVerification.Verification{}).Begin().
		Where("email = ?", email).
		Delete(&dbVerification.Verification{})
	err = query.Error
	if err != nil {
		query.Rollback()
		return
	}
	err = query.Commit().Error

	return
}

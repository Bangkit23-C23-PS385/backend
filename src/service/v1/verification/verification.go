package verification

import (
	"ta/backend/src/constant"
	dbVerification "ta/backend/src/entity/v1/db/verification"
	verifRepo "ta/backend/src/repository/v1/verification"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type Service struct {
	repo verifRepo.Repositorier
}

func NewService(repo verifRepo.Repositorier) *Service {
	return &Service{
		repo: repo,
	}
}

type Servicer interface {
	GetByEmail(email string) (verifEntity dbVerification.Verification, err error)
	Insert(email, token string) (err error)
	UpdateToken(email, token string, attemptLeft int) (err error)
	DeleteByEmail(email string) (err error)
}

func (svc Service) GetByEmail(email string) (verifEntity dbVerification.Verification, err error) {
	verifEntity, err = svc.repo.GetByEmail(email)
	if err == gorm.ErrRecordNotFound {
		err = constant.ErrTokenNotFound
		return
	} else if err != nil {
		err = errors.Wrap(err, "repo: get by username")
		return
	}

	return
}

func (svc Service) Insert(email, token string) (err error) {
	err = svc.repo.Insert(email, token)
	if err != nil {
		err = errors.Wrap(err, "repo: insert")
	}

	return
}

func (svc Service) UpdateToken(email, token string, attemptLeft int) (err error) {
	err = svc.repo.UpdateToken(email, token, attemptLeft)
	if err != nil {
		err = errors.Wrap(err, "repo: update token")
	}

	return
}

func (svc Service) DeleteByEmail(email string) (err error) {
	// check if exists
	_, err = svc.repo.GetByEmail(email)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = constant.ErrDataNotFound
		return
	} else if err != nil {
		err = errors.Wrap(err, "repo: get by email")
		return
	}

	// delete verification data
	err = svc.repo.DeleteByEmail(email)
	if err != nil {
		err = errors.Wrap(err, "delete by email")
	}

	return
}

package profile

import (
	"ta/backend/src/entity/v1/db/profile"
	// dbProfile "ta/backend/src/entity/v1/db/profile"
	profileRepo "ta/backend/src/repository/v1/profile"

	"github.com/pkg/errors"
)

type Service struct {
	// repo for profile
	repo profileRepo.Repositorier
}

func NewService(repo profileRepo.Repositorier) *Service {
	return &Service{
		repo: repo,
	}
}

type Servicer interface {
	GetProfile(id string) (*profile.Profile, error)
	CreateProfile(req profile.Profile) (*profile.Profile, error)
	UpdateProfile(req profile.Profile) (*profile.Profile, error)
	DeleteProfile(id string) (*profile.Profile, error)
}

// func (svc Service) ExtractToken(token string) (res profile.CommonRequest, err error) {
// }

func (svc Service) GetProfile(id string) (*profile.Profile, error) {
	var err error
	result := profile.Profile{}
	result, err = svc.repo.GetProfile(id)
	if err != nil {
		err = errors.Wrap(err, "repo: get profile")
		return nil, err
	}

	return &result, nil
}

func (svc Service) CreateProfile(req profile.Profile) (*profile.Profile, error) {
	var err error
	result, err := svc.GetProfile(req.UserId)
	if err != nil {
		err = errors.Wrap(err, "repo: get profile")
		return nil, err
	}
	if result != nil {
		err = errors.Wrap(err, "service: profile already exist")
		return nil, err
	}
	newProfile := profile.Profile{
		UserId:      req.UserId,
		Name:        req.Name,
		Gender:      req.Gender,
		DateOfBirth: req.DateOfBirth,
		Height:      req.Height,
		Weight:      req.Weight,
	}
	err = svc.repo.CreateProfile(newProfile)
	if err != nil {
		err = errors.Wrap(err, "repo: create profile")
		return nil, err
	}

	return &newProfile, nil
}

func (svc Service) UpdateProfile(req profile.Profile) (*profile.Profile, error) {
	var err error
	newProfile := profile.Profile{
		UserId:      req.UserId,
		Name:        req.Name,
		Gender:      req.Gender,
		DateOfBirth: req.DateOfBirth,
		Height:      req.Height,
		Weight:      req.Weight,
	}
	err = svc.repo.UpdateProfile(newProfile)
	if err != nil {
		err = errors.Wrap(err, "repo: update profile")
		return nil, err
	}

	return &newProfile, nil
}

func (svc Service) DeleteProfile(id string) (*profile.Profile, error) {
	var err error
	result, err := svc.GetProfile(id)
	if err != nil {
		err = errors.Wrap(err, "repo: get profile")
		return nil, err
	}
	if result == nil {
		err = errors.Wrap(err, "service: profile didnt exist")
		return nil, err
	}
	err = svc.repo.DeleteProfile(id)
	if err != nil {
		err = errors.Wrap(err, "repo: delete profile")
		return nil, err
	}

	return &profile.Profile{UserId: id}, nil
}

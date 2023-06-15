package profile

import (
	"backend/src/constant"
	"backend/src/entity/v1/db/profile"
	profileDto "backend/src/entity/v1/http/profile"
	"fmt"
	"time"

	// dbProfile "ta/backend/src/entity/v1/db/profile"
	profileRepo "backend/src/repository/v1/profile"

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
	GetProfile(id string) (*profileDto.ProfileDto, error)
	CreateProfile(req profile.Profile) (*profile.Profile, error)
	UpdateProfile(req profile.Profile) (*profile.Profile, error)
	DeleteProfile(id string) (*profileDto.ProfileDto, error)
	ConvertStringToTime(dateString string) (time.Time, error)
}

// func (svc Service) ExtractToken(token string) (res profile.CommonRequest, err error) {
// }

func (svc Service) GetProfile(id string) (*profileDto.ProfileDto, error) {
	var err error
	result := profile.Profile{}
	result, err = svc.repo.GetProfile(id)
	if err != nil {
		err = errors.Wrap(err, "repo: get profile")
		return nil, err
	}
	if result.CreatedAt.Equal(time.Time{}) {
		return nil, nil
	}
	response := &profileDto.ProfileDto{
		UserId:      result.UserId,
		Name:        result.Name,
		Gender:      result.Gender,
		DateOfBirth: result.DateOfBirth.Local().Format("2006-01-02"),
		Height:      result.Height,
		Weight:      result.Weight,
		CreatedAt:   result.CreatedAt,
		UpdatedAt:   result.UpdatedAt,
	}

	return response, nil
}

func (svc Service) CreateProfile(req profile.Profile) (*profile.Profile, error) {
	var err error
	result, err := svc.GetProfile(req.UserId)
	if err != nil {
		err = errors.Wrap(err, "repo: get profile")
		return nil, err
	}
	if result != nil {
		return nil, nil
	} else {
		newProfile := profile.Profile{}
		newProfile.UserId = req.UserId
		if req.Name != "" {
			newProfile.Name = req.Name
		}
		if req.Height != 0 {
			newProfile.Height = req.Height
		}
		if req.Weight != 0 {
			newProfile.Weight = req.Weight
		}
		if req.Gender != "" {
			newProfile.Gender = req.Gender
		}
		if !req.DateOfBirth.Equal(time.Time{}) {
			newProfile.DateOfBirth = req.DateOfBirth
		}
		err = svc.repo.CreateProfile(newProfile)
		if err != nil {
			err = errors.Wrap(err, "repo: create profile")
			return nil, err
		}

		return &newProfile, nil
	}

}

func (svc Service) UpdateProfile(req profile.Profile) (*profile.Profile, error) {
	var err error
	result, err := svc.GetProfile(req.UserId)
	if err != nil {
		err = errors.Wrap(err, "repo: get profile")
		return nil, err
	}
	if result == nil {
		return nil, nil
	} else {
		// Parse the string into a time.Time value
		t, err := time.Parse("2006-01-02", result.DateOfBirth)
		if err != nil {
			fmt.Println("Error parsing date:", err)
			err = errors.Wrap(err, "Error parsing date:")
			return nil, err
		}
		newProfile := &profile.Profile{
			UserId:      result.UserId,
			Name:        result.Name,
			Gender:      result.Gender,
			DateOfBirth: t,
			Height:      result.Height,
			Weight:      result.Weight,
			CreatedAt:   result.CreatedAt,
			UpdatedAt:   result.UpdatedAt,
		}
		if req.Name != "" {
			newProfile.Name = req.Name
		}
		if !req.DateOfBirth.Equal(time.Time{}) {
			newProfile.DateOfBirth = req.DateOfBirth
		}

		if req.Height != 0 {
			newProfile.Height = req.Height
		}
		if req.Weight != 0 {
			newProfile.Weight = req.Weight
		}
		if req.Gender != "" {
			if req.Gender != constant.Lakilaki && req.Gender != constant.Perempuan {
				err = errors.Wrap(err, "wrong gender value")
				return nil, err
			}
			newProfile.Gender = req.Gender
		}
		newProfile.UpdatedAt = time.Now()

		err = svc.repo.UpdateProfile(*newProfile)
		if err != nil {
			err = errors.Wrap(err, "repo: update profile")
			return nil, err
		}

		return newProfile, nil
	}

}

func (svc Service) DeleteProfile(id string) (*profileDto.ProfileDto, error) {
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

	return &profileDto.ProfileDto{
		UserId:      id,
		Name:        result.Name,
		Gender:      result.Gender,
		DateOfBirth: result.DateOfBirth,
		Height:      result.Height,
		Weight:      result.Weight,
		CreatedAt:   result.CreatedAt,
		UpdatedAt:   result.UpdatedAt,
	}, nil
}

func (svc Service) ConvertStringToTime(dateString string) (time.Time, error) {
	layout := "2006-01-02" // Format for "YYYY-MM-DD"
	t, err := time.Parse(layout, dateString)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}

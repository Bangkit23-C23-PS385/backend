package profile

import (
	"backend/src/constant"
)

type GetDeleteCommonRequest struct {
	UserID string `json:"id"`
}

type CreateRequest struct {
	Gender      constant.GenderType `json:"gender" example:"LAKILAKI/PEREMPUAN"`
	DateOfBirth string              `json:"dateOfBirth" example:"YYYY-MM-DD"`
	Height      int                 `json:"height"`
	Weight      int                 `json:"weight"`
}
type UpdateRequest struct {
	Name        string              `json:"name"`
	Gender      constant.GenderType `json:"gender" example:""`
	DateOfBirth string              `json:"dateOfBirth" example:""`
	Height      int                 `json:"height"`
	Weight      int                 `json:"weight"`
}

// making type Code for inside CommonResponse
type Code string

const (
	Success             Code = "SUCCESS"
	ProfileAlreadyExist Code = "PROFILE_ALREADY_EXIST"
	BadRequest          Code = "BAD_REQUEST"
	NotFound            Code = "NOT_FOUND"
	SystemError         Code = "SYSTEM_ERROR"
)

// CommonResponse is the response structure for all the CRUD operations
type CommonResponse struct {
	Code    Code        `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

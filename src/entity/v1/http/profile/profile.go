package profile

import (
	"ta/backend/src/constant"
	"time"
)

type GetDeleteCommonRequest struct {
	UserID string `json:"id"`
}

type CommonRequest struct {
	UserID      string              `json:"id"`
	Name        string              `json:"column:name"`
	Gender      constant.GenderType `json:"column:gender"`
	DateOfBirth time.Time           `json:"column:dateOfBirth"`
	Height      int                 `json:"column:height"`
	Weight      string              `json:"column:weight"`
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

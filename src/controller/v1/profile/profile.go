package profile

import (
	"backend/src/constant"
	profileModel "backend/src/entity/v1/db/profile"
	profileDto "backend/src/entity/v1/http/profile"
	"backend/src/service/v1/profile"
	"backend/src/util/helper"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/pkg/errors"
)

type Controller struct {
	svc profile.Servicer
}

func NewController(svc profile.Servicer) *Controller {
	return &Controller{
		svc: svc,
	}
}

// Get Profile godoc
// @Summary Get Profile
// @Description Get Profile
// @Tags Profile
// @Param Authorization header string true "Bearer Token"
// @Success 200 {string} string "OK"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Router /v1/profile [get]
func (ctrl Controller) GetProfile(ctx *gin.Context) {
	token, err := helper.GetJWT(ctx)
	if errors.Is(err, constant.ErrTokenNotFound) || errors.Is(err, constant.ErrInvalidFormat) {
		log.Println(err)
		helper.JSONResponse(ctx, http.StatusUnauthorized, errors.Cause(err).Error(), nil)
		return
	}
	claims, err := helper.ExtractJWT(token)
	if errors.Is(err, jwt.ErrSignatureInvalid) || errors.Is(err, constant.ErrTokenInvalid) {
		log.Println(err)
		helper.JSONResponse(ctx, http.StatusUnauthorized, errors.Cause(err).Error(), nil)
		return
	}

	res, err := ctrl.svc.GetProfile(claims.UserID)
	log.Println(res.Name)
	if err != nil {
		log.Println(err)
		helper.JSONResponse(ctx, http.StatusInternalServerError, errors.Cause(err).Error(), nil)
	} else {
		log.Println(res.Name)
		helper.JSONResponse(ctx, http.StatusOK, "", res)
	}
}

// Create Profile godoc
// @Summary Create Profile
// @Description Create Profile
// @Tags Profile
// @Param Authorization header string true "Bearer Token"
// @Param profile body CreateRequest true "Profile request body"
// @Success 200 {string} string "OK"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Router /v1/profile [post]
func (ctrl Controller) CreateProfile(ctx *gin.Context) {
	token, err := helper.GetJWT(ctx)
	if errors.Is(err, constant.ErrTokenNotFound) || errors.Is(err, constant.ErrInvalidFormat) {
		log.Println(err)
		helper.JSONResponse(ctx, http.StatusUnauthorized, errors.Cause(err).Error(), nil)
		return
	}
	claims, err := helper.ExtractJWT(token)
	if errors.Is(err, jwt.ErrSignatureInvalid) || errors.Is(err, constant.ErrTokenInvalid) {
		log.Println(err)
		helper.JSONResponse(ctx, http.StatusUnauthorized, errors.Cause(err).Error(), nil)
		return
	}
	var request profileDto.CreateRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		log.Println(err)
		helper.JSONResponse(ctx, http.StatusBadRequest, "Invalid request body", nil)
		return
	}
	dateTime := time.Time{}
	if request.DateOfBirth == "YYYY-MM-DD" || request.DateOfBirth == "" {
		dateTime = time.Time{}
	} else {
		dateTime, err = ctrl.svc.ConvertStringToTime(request.DateOfBirth)
		if err != nil {
			log.Println(err)
			helper.JSONResponse(ctx, http.StatusInternalServerError, "failed to parse datestring to datetime", nil)
			return
		}
	}
	res, err := ctrl.svc.CreateProfile(profileModel.Profile{
		UserId:      claims.UserID,
		Name:        claims.UserName,
		Gender:      request.Gender,
		DateOfBirth: dateTime,
		Height:      request.Height,
		Weight:      request.Weight,
	})
	if err != nil {
		log.Println(err)
		helper.JSONResponse(ctx, http.StatusInternalServerError, errors.Cause(err).Error(), nil)
	}
	if err == nil && res != nil {
		log.Println(res.Name)
		helper.JSONResponse(ctx, http.StatusOK, "", res)
	}
	if err == nil && res == nil {
		helper.JSONResponse(ctx, http.StatusBadRequest, errors.New("Profile Already Exist").Error(), nil)
	}
}

// Update Profile godoc
// @Summary Update Profile
// @Description Update Profile
// @Tags Profile
// @Param Authorization header string true "Bearer Token"
// @Param profile body UpdateRequest true "Profile request body"
// @Success 200 {string} string "OK"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Router /v1/profile [put]
func (ctrl Controller) UpdateProfile(ctx *gin.Context) {
	token, err := helper.GetJWT(ctx)
	if errors.Is(err, constant.ErrTokenNotFound) || errors.Is(err, constant.ErrInvalidFormat) {
		log.Println(err)
		helper.JSONResponse(ctx, http.StatusUnauthorized, errors.Cause(err).Error(), nil)
		return
	}
	claims, err := helper.ExtractJWT(token)
	if errors.Is(err, jwt.ErrSignatureInvalid) || errors.Is(err, constant.ErrTokenInvalid) {
		log.Println(err)
		helper.JSONResponse(ctx, http.StatusUnauthorized, errors.Cause(err).Error(), nil)
		return
	}
	var request profileDto.UpdateRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		log.Println(err)
		helper.JSONResponse(ctx, http.StatusBadRequest, "Invalid request body", nil)
		return
	}
	dateTime := time.Time{}
	if request.DateOfBirth == "YYYY-MM-DD" || request.DateOfBirth == "" {
		dateTime = time.Time{}
	} else {
		dateTime, err = ctrl.svc.ConvertStringToTime(request.DateOfBirth)
		if err != nil {
			log.Println(err)
			helper.JSONResponse(ctx, http.StatusInternalServerError, "failed to parse datestring to datetime", nil)
			return
		}
	}

	res, err := ctrl.svc.UpdateProfile(profileModel.Profile{
		UserId:      claims.UserID,
		Name:        request.Name,
		Gender:      request.Gender,
		DateOfBirth: dateTime,
		Height:      request.Height,
		Weight:      request.Weight,
	})
	if err != nil {
		log.Println(err)
		helper.JSONResponse(ctx, http.StatusInternalServerError, errors.Cause(err).Error(), nil)
	}
	if err == nil && res != nil {
		log.Println(res.Name)
		helper.JSONResponse(ctx, http.StatusOK, "", res)
	}
	if err == nil && res == nil {
		helper.JSONResponse(ctx, http.StatusBadRequest, errors.New("Profile Didnt Exist").Error(), nil)
	}
}

// Delete Profile godoc
// @Summary Delete Profile
// @Description Delete Profile
// @Tags Profile
// @Param Authorization header string true "Bearer Token"
// @Success 200 {string} string "OK"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Router /v1/profile [delete]
func (ctrl Controller) DeleteProfile(ctx *gin.Context) {
	token, err := helper.GetJWT(ctx)
	if errors.Is(err, constant.ErrTokenNotFound) || errors.Is(err, constant.ErrInvalidFormat) {
		log.Println(err)
		helper.JSONResponse(ctx, http.StatusUnauthorized, errors.Cause(err).Error(), nil)
		return
	}
	claims, err := helper.ExtractJWT(token)
	if errors.Is(err, jwt.ErrSignatureInvalid) || errors.Is(err, constant.ErrTokenInvalid) {
		log.Println(err)
		helper.JSONResponse(ctx, http.StatusUnauthorized, errors.Cause(err).Error(), nil)
		return
	}

	res, err := ctrl.svc.DeleteProfile(claims.UserID)
	if err != nil {
		log.Println(err)
		helper.JSONResponse(ctx, http.StatusInternalServerError, errors.Cause(err).Error(), nil)
	} else {
		log.Println(res.Name)
		helper.JSONResponse(ctx, http.StatusOK, "", res)
	}
}

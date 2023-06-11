package predict

import (
	"backend/src/constant"
	httpPredict "backend/src/entity/v1/http/predict"
	predictSvc "backend/src/service/v1/predict"
	"backend/src/util/helper"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
)

type Controller struct {
	svc predictSvc.Servicer
}

func NewController(svc predictSvc.Servicer) *Controller {
	return &Controller{
		svc: svc,
	}
}

// Get Symptoms godoc
// @Summary Get Symptoms
// @Description Get Symptoms
// @Tags Predict
// @Param Authorization header string true "Bearer Token"
// @Success 200 {string} string "OK"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Router /v1/symptoms [get]
func (ctrl Controller) GetSymptoms(ctx *gin.Context) {
	token, err := helper.GetJWT(ctx)
	if errors.Is(err, constant.ErrTokenNotFound) || errors.Is(err, constant.ErrInvalidFormat) {
		log.Println(err)
		helper.JSONResponse(ctx, http.StatusUnauthorized, errors.Cause(err).Error(), nil)
		return
	}
	_, err = helper.ExtractJWT(token)
	if errors.Is(err, jwt.ErrSignatureInvalid) || errors.Is(err, constant.ErrTokenInvalid) {
		log.Println(err)
		helper.JSONResponse(ctx, http.StatusUnauthorized, errors.Cause(err).Error(), nil)
		return
	}

	resps, err := ctrl.svc.GetSymptoms()
	if err != nil {
		log.Println(err)
		helper.JSONResponse(ctx, http.StatusInternalServerError, errors.Cause(err).Error(), nil)
	} else {
		helper.JSONResponse(ctx, http.StatusOK, "", resps)
	}
}

// SubmitData godoc
// @Summary Submit Data
// @Description Submit Data
// @Tags Predict
// @Param Authorization header string true "Bearer Token"
// @Param Payload body predict.PredictSymptoms true "Payload"
// @Success 200 {string} string "OK"
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Router /v1/predict [post]
func (ctrl Controller) SubmitData(ctx *gin.Context) {
	token, err := helper.GetJWT(ctx)
	if errors.Is(err, constant.ErrTokenNotFound) || errors.Is(err, constant.ErrInvalidFormat) {
		log.Println(err)
		helper.JSONResponse(ctx, http.StatusUnauthorized, errors.Cause(err).Error(), nil)
		return
	}
	_, err = helper.ExtractJWT(token)
	if errors.Is(err, jwt.ErrSignatureInvalid) || errors.Is(err, constant.ErrTokenInvalid) {
		log.Println(err)
		helper.JSONResponse(ctx, http.StatusUnauthorized, errors.Cause(err).Error(), nil)
		return
	}

	req := httpPredict.PredictSymptoms{}
	err = ctx.BindJSON(&req)
	if err != nil {
		log.Println(err)
		helper.JSONResponse(ctx, http.StatusBadRequest, constant.ErrInvalidFormat.Error(), nil)
		return
	}
	if len(req.Symptoms) < 3 || len(req.Symptoms) > 7 {
		helper.JSONResponse(ctx, http.StatusBadRequest, constant.ErrInvalidSymptomLength.Error(), nil)
		return
	}

	err = ctrl.svc.SubmitData(req)
	if err != nil {
		log.Println(err)
		helper.JSONResponse(ctx, http.StatusInternalServerError, errors.Cause(err).Error(), nil)
	} else {
		helper.JSONResponse(ctx, http.StatusOK, "", nil)
	}
}

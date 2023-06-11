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

// SubmitData godoc
// @Summary Submit Rating
// @Description Submit candidate rating with its interview result if it's their last interview
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

	err = ctrl.svc.SubmitData(req)
	if err != nil {
		log.Println(err)
		helper.JSONResponse(ctx, http.StatusInternalServerError, errors.Cause(err).Error(), nil)
	} else {
		helper.JSONResponse(ctx, http.StatusOK, "", nil)
	}
}

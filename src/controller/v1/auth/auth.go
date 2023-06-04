package auth

import (
	"backend/src/constant"
	httpAuth "backend/src/entity/v1/http/auth"
	authSvc "backend/src/service/v1/auth"
	"backend/src/util/helper"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
)

type Controller struct {
	svc authSvc.Servicer
}

func NewController(
	svc authSvc.Servicer,
) *Controller {
	return &Controller{
		svc: svc,
	}
}

// Login godoc
// @Summary Login
// @Description Login
// @Tags Auth
// @Param Payload body auth.LoginRequest true "Payload"
// @Success 200 {object} auth.LoginResponse "OK"
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Router /v1/login [post]
func (ctrl Controller) Login(ctx *gin.Context) {
	req := httpAuth.LoginRequest{}
	err := ctx.BindJSON(&req)
	if err != nil {
		log.Println(err)
		helper.JSONResponse(ctx, http.StatusBadRequest, constant.ErrInvalidFormat.Error(), nil)
		return
	}

	val := validator.New()
	err = val.Struct(req)
	if err != nil {
		log.Println(err)
		helper.JSONResponse(ctx, http.StatusBadRequest, constant.ErrInvalidFormat.Error(), nil)
		return
	}

	if helper.IsIdentifierEmail(req.Identifier) {
		err = helper.SanitizeEmail(req.Identifier)
		if errors.Is(err, constant.ErrEmailLength) || errors.Is(err, constant.ErrEmailInvalid) {
			log.Println(err)
			helper.JSONResponse(ctx, http.StatusBadRequest, errors.Cause(err).Error(), nil)
			return
		} else if err != nil {
			log.Println(err)
			helper.JSONResponse(ctx, http.StatusInternalServerError, errors.Cause(err).Error(), nil)
		}
	} else {
		err = helper.SanitizeUsername(req.Identifier)
		if errors.Is(err, constant.ErrUsernameLength) || errors.Is(err, constant.ErrUsernameUnallowed) ||
			errors.Is(err, constant.ErrUsernamePrefix) || errors.Is(err, constant.ErrUsernameSuffix) {
			log.Println(err)
			helper.JSONResponse(ctx, http.StatusBadRequest, errors.Cause(err).Error(), nil)
			return
		} else if err != nil {
			log.Println(err)
			helper.JSONResponse(ctx, http.StatusInternalServerError, errors.Cause(err).Error(), nil)
			return
		}
	}

	loginResponse, err := ctrl.svc.Login(req)
	if errors.Is(err, constant.ErrAccountNotVerified) {
		log.Println(err)
		helper.JSONResponse(ctx, http.StatusUnauthorized, constant.ErrAccountNotVerified.Error(), nil)
	} else if errors.Is(err, constant.ErrInvalidCreds) || errors.Is(err, constant.ErrAccountNotFound) {
		log.Println(err)
		helper.JSONResponse(ctx, http.StatusBadRequest, constant.ErrInvalidCreds.Error(), nil)
	} else if err != nil {
		log.Println(err)
		helper.JSONResponse(ctx, http.StatusInternalServerError, errors.Cause(err).Error(), nil)
	} else {
		ctx.SetCookie("Token", "", 0, "/", constant.ApplicationDomain, false, true)
		ctx.SetCookie("Token", loginResponse.AccessToken, constant.AccessTokenInterval, "/", constant.ApplicationDomain, false, true)
		helper.JSONResponse(ctx, http.StatusOK, "", loginResponse)
	}
}

// Logout godoc
// @Summary Logout
// @Description Logout
// @Tags Auth
// @Param Authorization header string true "Bearer Token"
// @Success 200 {string} string "OK"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Router /v1/logout [post]
func (ctrl Controller) Logout(ctx *gin.Context) {
	claims, err := helper.GetClaims(ctx)
	if errors.Is(err, constant.ErrTokenNotFound) || errors.Is(err, constant.ErrTokenInvalid) ||
		errors.Is(err, constant.ErrTokenNotFound) || errors.Is(err, constant.ErrInvalidFormat) {
		log.Println(err)
		helper.JSONResponse(ctx, http.StatusUnauthorized, errors.Cause(err).Error(), nil)
		return
	} else if err != nil {
		log.Println(err)
		helper.JSONResponse(ctx, http.StatusInternalServerError, errors.Cause(err).Error(), nil)
		return
	}

	err = ctrl.svc.Logout(claims)
	if err != nil {
		log.Println(err)
		helper.JSONResponse(ctx, http.StatusInternalServerError, errors.Cause(err).Error(), nil)
	} else {
		ctx.SetCookie("Token", "", 0, "/", constant.ApplicationDomain, false, true)
		helper.JSONResponse(ctx, http.StatusOK, "", nil)
	}
}

// Register godoc
// @Summary Register
// @Description Register
// @Tags Auth
// @Param Payload body auth.RegisterRequest true "Payload"
// @Success 201 {string} string "Created"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /v1/register [post]
func (ctrl Controller) Register(ctx *gin.Context) {
	req := httpAuth.RegisterRequest{}
	err := ctx.BindJSON(&req)
	if err != nil {
		log.Println(err)
		helper.JSONResponse(ctx, http.StatusBadRequest, constant.ErrInvalidFormat.Error(), nil)
		return
	}

	val := validator.New()
	err = val.Struct(req)
	if err != nil {
		log.Println(err)
		helper.JSONResponse(ctx, http.StatusBadRequest, constant.ErrInvalidFormat.Error(), nil)
		return
	}

	err = helper.SanitizeEmail(req.Email)
	if errors.Is(err, constant.ErrEmailLength) || errors.Is(err, constant.ErrEmailInvalid) {
		log.Println(err)
		helper.JSONResponse(ctx, http.StatusBadRequest, errors.Cause(err).Error(), nil)
		return
	} else if err != nil {
		log.Println(err)
		helper.JSONResponse(ctx, http.StatusInternalServerError, errors.Cause(err).Error(), nil)
		return
	}

	err = helper.SanitizeUsername(req.Username)
	if errors.Is(err, constant.ErrUsernameLength) || errors.Is(err, constant.ErrUsernameUnallowed) ||
		errors.Is(err, constant.ErrUsernamePrefix) || errors.Is(err, constant.ErrUsernameSuffix) {
		log.Println(err)
		helper.JSONResponse(ctx, http.StatusBadRequest, errors.Cause(err).Error(), nil)
		return
	} else if err != nil {
		log.Println(err)
		helper.JSONResponse(ctx, http.StatusInternalServerError, errors.Cause(err).Error(), nil)
		return
	}

	err = helper.SanitizePassword(req.Password)
	if errors.Is(err, constant.ErrPasswordLength) || errors.Is(err, constant.ErrPasswordUnallowed) ||
		errors.Is(err, constant.ErrPasswordWeak) {
		log.Println(err)
		helper.JSONResponse(ctx, http.StatusBadRequest, errors.Cause(err).Error(), nil)
		return
	} else if err != nil {
		log.Println(err)
		helper.JSONResponse(ctx, http.StatusInternalServerError, errors.Cause(err).Error(), nil)
		return
	}

	err = ctrl.svc.Register(req)
	if errors.Is(err, constant.ErrEmailTaken) {
		log.Println(err)
		helper.JSONResponse(ctx, http.StatusBadRequest, errors.Cause(err).Error(), nil)
	} else if err != nil {
		log.Println(err)
		helper.JSONResponse(ctx, http.StatusInternalServerError, errors.Cause(err).Error(), nil)
	} else {
		helper.JSONResponse(ctx, http.StatusCreated, "", nil)
	}
}

// Verify Email godoc
// @Summary Verify Email
// @Description Verify Email
// @Tags Auth
// @Param username query string true "Username"
// @Param token query string true "Token"
// @Success 302 {string} string "Found"
// @Failure 400 {string} string "Bad Request"
// @Failure 422 {string} string "Unprocessable Entity"
// @Failure 500 {string} string "Internal Server Error"
// @Router /v1/verify [get]
func (ctrl Controller) VerifyEmail(ctx *gin.Context) {
	req := httpAuth.VerifyRequest{}
	err := ctx.BindQuery(&req)
	if err != nil {
		log.Print(err)
		helper.JSONResponse(ctx, http.StatusBadRequest, constant.ErrInvalidFormat.Error(), nil)
		return
	}

	val := validator.New()
	err = val.Struct(req)
	if err != nil {
		log.Println(err)
		helper.JSONResponse(ctx, http.StatusBadRequest, constant.ErrInvalidFormat.Error(), nil)
		return
	}

	err = ctrl.svc.VerifyEmail(req)
	if errors.Is(err, constant.ErrAccountNotFound) || errors.Is(err, constant.ErrAccountAlreadyVerified) ||
		errors.Is(err, constant.ErrTokenInvalid) || errors.Is(err, constant.ErrTokenNotFound) {
		log.Println(err)
		helper.JSONResponse(ctx, http.StatusUnprocessableEntity, errors.Cause(err).Error(), nil)
	} else if err != nil {
		log.Println(err)
		helper.JSONResponse(ctx, http.StatusInternalServerError, errors.Cause(err).Error(), nil)
	} else {
		// location := fmt.Sprintf("%v", constant.RedirectURL)
		// ctx.Redirect(http.StatusFound, location)
		helper.JSONResponse(ctx, http.StatusOK, "Verifikasi berhasil. Silakan login ke aplikasi", nil)
	}
}

// Resend godoc
// @Summary Resend Email
// @Description Resend email for HR verification
// @Tags Auth
// @Param Payload body auth.ResendRequest true "Payload"
// @Success 200 {string} string "OK"
// @Failure 400 {string} string "Bad Request"
// @Failure 403 {string} string "Forbidden"
// @Failure 422 {string} string "Unprocessable Entity"
// @Failure 500 {string} string "Internal Server Error"
// @Router /v1/resend [post]
func (ctrl Controller) Resend(ctx *gin.Context) {
	req := httpAuth.ResendRequest{}
	err := ctx.BindJSON(&req)
	if err != nil {
		log.Print(err)
		helper.JSONResponse(ctx, http.StatusBadRequest, constant.ErrInvalidFormat.Error(), nil)
		return
	}

	val := validator.New()
	err = val.Struct(req)
	if err != nil {
		log.Println(err)
		helper.JSONResponse(ctx, http.StatusBadRequest, constant.ErrInvalidFormat.Error(), nil)
		return
	}

	err = ctrl.svc.Resend(req)
	if errors.Is(err, constant.ErrAccountNotFound) || errors.Is(err, constant.ErrAccountAlreadyVerified) {
		log.Println(err)
		helper.JSONResponse(ctx, http.StatusUnprocessableEntity, errors.Cause(err).Error(), nil)
	} else if errors.Is(err, constant.ErrNoResendAttemptLeft) {
		log.Println(err)
		helper.JSONResponse(ctx, http.StatusForbidden, errors.Cause(err).Error(), nil)
	} else if err != nil {
		log.Println(err)
		helper.JSONResponse(ctx, http.StatusInternalServerError, errors.Cause(err).Error(), nil)
	} else {
		helper.JSONResponse(ctx, http.StatusOK, "", nil)
	}
}

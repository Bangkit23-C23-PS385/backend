package auth

import (
	"backend/src/constant"
	dbUser "backend/src/entity/v1/db/user"
	httpAuth "backend/src/entity/v1/http/auth"
	authRepo "backend/src/repository/v1/auth"
	verifSvc "backend/src/service/v1/verification"
	"backend/src/util/email"
	"backend/src/util/helper"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type Service struct {
	verifSvc verifSvc.Servicer
	repo     authRepo.Repositorier
}

func NewService(
	verifSvc verifSvc.Servicer,
	repo authRepo.Repositorier,
) *Service {
	return &Service{
		verifSvc: verifSvc,
		repo:     repo,
	}
}

type Servicer interface {
	Login(req httpAuth.LoginRequest) (resp httpAuth.LoginResponse, err error)
	Logout(claims httpAuth.Claims) (err error)
	Register(req httpAuth.RegisterRequest) (err error)
	VerifyEmail(req httpAuth.VerifyRequest) (err error)
	Resend(req httpAuth.ResendRequest) (err error)
}

func (svc Service) Login(req httpAuth.LoginRequest) (resp httpAuth.LoginResponse, err error) {
	// get all credentials by email
	var userEntity dbUser.User
	if helper.IsIdentifierEmail(req.Identifier) {
		userEntity, err = svc.repo.GetUserByEmail(req.Identifier)
	} else {
		userEntity, err = svc.repo.GetUserByUsername(req.Identifier)
	}
	if err == gorm.ErrRecordNotFound {
		err = constant.ErrAccountNotFound
		return
	} else if !userEntity.IsVerified {
		err = constant.ErrAccountNotVerified
		return
	} else if err != nil {
		err = errors.Wrap(err, "repo: get user by email or username")
		return
	}

	// validate password
	hashedPassword := helper.HashPassword(req.Password, userEntity.PasswordSalt)
	if hashedPassword != userEntity.Password {
		err = constant.ErrInvalidCreds
		return
	}

	// generate token
	encryptedID, err := helper.Encrypt(strconv.Itoa(int(userEntity.ID)))
	if err != nil {
		err = errors.Wrap(err, "aes: encrypt")
		return
	}
	expiresAt := time.Now().Add(constant.AccessTokenDuration)
	access_token, err := svc.generateToken(encryptedID, userEntity.Name, userEntity.Email, userEntity.Username, expiresAt)
	if err != nil {
		err = errors.Wrap(err, "svc: generate token")
		return
	}

	resp = httpAuth.LoginResponse{
		UserResponse: httpAuth.UserResponse{
			ID:       encryptedID,
			Name:     userEntity.Name,
			Email:    userEntity.Email,
			Username: userEntity.Username,
		},
		AccessToken: access_token,
	}

	return

}

func (svc Service) generateToken(encryptedID, name, email, username string, expiresAt time.Time) (token string, err error) {
	issuer := constant.ApplicationName
	claim := &httpAuth.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    issuer,
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
		UserID:       encryptedID,
		UserName:     name,
		UserEmail:    email,
		UserUsername: username,
	}

	generatedToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	token, err = generatedToken.SignedString([]byte(os.Getenv("JWT_ACCESS_SIGNATURE_KEY")))
	if err != nil {
		err = errors.Wrap(err, "jwt: generating signed token")
	}

	return
}

func (svc Service) Logout(claims httpAuth.Claims) (err error) {
	expiresAt := time.Now()
	_, err = svc.generateToken(claims.ID, claims.UserName, claims.UserEmail, claims.UserUsername, expiresAt)
	if err != nil {
		err = errors.Wrap(err, "svc: generate token")
	}

	return
}

func (svc Service) Register(req httpAuth.RegisterRequest) (err error) {
	// check username and email availability
	userEntity, err := svc.repo.GetUserByEmail(req.Email)
	if err == gorm.ErrRecordNotFound {
		err = nil
	} else if userEntity.Email == req.Email {
		err = constant.ErrEmailTaken
		return
	} else if err != nil {
		err = errors.Wrap(err, "repo: get user by email")
		return
	}
	userEntity, err = svc.repo.GetUserByUsername(req.Username)
	if err == gorm.ErrRecordNotFound {
		err = nil
	} else if userEntity.Username == req.Username {
		err = constant.ErrUsernameTaken
		return
	} else if err != nil {
		err = errors.Wrap(err, "repo: get user by username")
		return
	}

	// salt password
	req.PasswordSalt = helper.GenerateSalt()
	req.Password = helper.HashPassword(req.Password, req.PasswordSalt)

	// insert to database
	err = svc.repo.Insert(req)
	if err != nil {
		err = errors.Wrap(err, "repo: insert user")
		return
	}

	// handle token generation for email verif
	verifToken := helper.GenerateEmailVerifToken()
	err = svc.verifSvc.Insert(req.Email, verifToken)
	if err != nil {
		err = errors.Wrap(err, "verif svc: insert")
		return
	}

	// send verification email to user
	link := fmt.Sprintf("https://%v/v1/verify?username=%v&token=%v", constant.APIOriginURL, req.Username, verifToken)
	data := map[string]interface{}{
		"name": req.Username,
		"link": link,
	}
	emailContent, err := email.LoadTemplate(string(constant.Register), data)
	if err != nil {
		err = errors.Wrap(err, "email: load template")
		return
	}
	err = email.SendEmail(req.Email, string(constant.VerifyEmail), emailContent)
	if err != nil {
		err = errors.Wrap(err, "email: send email")
	}

	return
}

func (svc Service) VerifyEmail(req httpAuth.VerifyRequest) (err error) {
	// get user data from db
	userEntity, err := svc.repo.GetUserByUsername(req.Username)
	if err == gorm.ErrRecordNotFound {
		err = constant.ErrAccountNotFound
		return
	} else if userEntity.IsVerified {
		err = constant.ErrAccountAlreadyVerified
		return
	} else if err != nil {
		err = errors.Wrap(err, "repo: get user by username")
		return
	}

	// get verification data from db
	verifEntity, err := svc.verifSvc.GetByEmail(userEntity.Email)
	if err != nil {
		err = errors.Wrap(err, "verif svc: get by email")
		return
	}

	// verify token
	if req.Token != verifEntity.Token {
		err = constant.ErrTokenInvalid
		return
	}

	userEntity.IsVerified = true

	// update user's verification status
	err = svc.repo.Update(userEntity)
	if err != nil {
		err = errors.Wrap(err, "repo: verify user")
		return
	}

	return
}

func (svc Service) Resend(req httpAuth.ResendRequest) (err error) {
	// get user data from db
	userEntity, err := svc.repo.GetUserByEmail(req.Email)
	if err == gorm.ErrRecordNotFound {
		err = constant.ErrAccountNotFound
		return
	} else if userEntity.IsVerified {
		err = constant.ErrAccountAlreadyVerified
		return
	} else if err != nil {
		err = errors.Wrap(err, "repo: get user by email")
		return
	}

	// get verification data from db
	verifEntity, err := svc.verifSvc.GetByEmail(userEntity.Email)
	if verifEntity.AttemptLeft <= 0 {
		err = constant.ErrNoResendAttemptLeft
		return
	} else if err != nil {
		err = errors.Wrap(err, "verif svc: get by email")
		return
	}

	// handle token generation for resend email verif
	verifToken := helper.GenerateEmailVerifToken()
	err = svc.verifSvc.UpdateToken(req.Email, verifToken, verifEntity.AttemptLeft-1)
	if err != nil {
		err = errors.Wrap(err, "verif svc: insert")
		return
	}

	// send verification email to user
	link := fmt.Sprintf("http://%v/v1/verify?username=%v&token=%v", constant.APIOriginURL, userEntity.Username, verifToken)
	data := map[string]interface{}{
		"name": userEntity.Username,
		"link": link,
	}
	emailContent, err := email.LoadTemplate(string(constant.Register), data)
	if err != nil {
		err = errors.Wrap(err, "email: load template")
		return
	}
	err = email.SendEmail(userEntity.Email, string(constant.VerifyEmail), emailContent)
	if err != nil {
		err = errors.Wrap(err, "email: send email")
	}

	return
}

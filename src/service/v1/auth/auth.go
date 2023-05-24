package auth

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"ta/backend/src/constant"
	dbUser "ta/backend/src/entity/v1/db/user"
	httpAuth "ta/backend/src/entity/v1/http/auth"
	authRepo "ta/backend/src/repository/v1/auth"
	verifSvc "ta/backend/src/service/v1/verification"
	"ta/backend/src/util/email"
	"ta/backend/src/util/helper"
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
	userEntity, err := svc.repo.GetUserByEmail(req.Email)
	if err == gorm.ErrRecordNotFound {
		err = constant.ErrAccountNotFound
		return
	} else if !userEntity.IsVerified {
		err = constant.ErrAccountNotVerified
		return
	} else if err != nil {
		err = errors.Wrap(err, "repo: get user by email")
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
	access_token, err := svc.generateToken(encryptedID, userEntity.Name, userEntity.Email, userEntity.Role, expiresAt)
	if err != nil {
		err = errors.Wrap(err, "svc: generate token")
		return
	}

	resp = httpAuth.LoginResponse{
		UserResponse: httpAuth.UserResponse{
			ID:    encryptedID,
			Name:  userEntity.Name,
			Email: userEntity.Email,
			Role:  userEntity.Role.String(),
		},
		AccessToken: access_token,
	}

	return

}

func (svc Service) generateToken(encryptedID, name, email string, role constant.Roles, expiresAt time.Time) (token string, err error) {
	issuer := constant.ApplicationName
	claim := &httpAuth.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    issuer,
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
		UserID:    encryptedID,
		UserName:  name,
		UserEmail: email,
		UserRole:  role.String(),
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
	_, err = svc.generateToken(claims.ID, claims.UserName, claims.UserEmail, constant.Roles(claims.UserRole), expiresAt)
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
	link := fmt.Sprintf("http://%v/v1/verify?email=%v&token=%v", constant.APIOriginURL, req.Email, verifToken)
	splitEmail := strings.Split(req.Email, "@")
	data := map[string]interface{}{
		"name": splitEmail[0],
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
	if err != nil {
		err = errors.Wrap(err, "verif svc: get by email")
		return
	}

	// verify token
	if req.Token != verifEntity.Token {
		err = constant.ErrTokenInvalid
		return
	}

	// generate password and hash it
	password := helper.GenereateRandomString()
	userEntity.PasswordSalt = helper.GenerateSalt()
	userEntity.Password = helper.HashPassword(password, userEntity.PasswordSalt)
	userEntity.IsVerified = true

	// update user's verification status
	err = svc.repo.Update(userEntity)
	if err != nil {
		err = errors.Wrap(err, "repo: verify user")
		return
	}

	// send email containing user's newly generated password
	emailSplit := strings.Split(userEntity.Email, "@")
	data := map[string]interface{}{
		"name":     emailSplit[0],
		"password": password,
	}
	emailContent, err := email.LoadTemplate(string(constant.VerifySuccess), data)
	if err != nil {
		err = errors.Wrap(err, "email: load template")
		return
	}
	err = email.SendEmail(req.Email, string(constant.RegistrationSuccess), emailContent)
	if err != nil {
		err = errors.Wrap(err, "email: send email")
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
	link := fmt.Sprintf("http://%v/v1/verify?email=%v&token=%v", constant.APIOriginURL, req.Email, verifToken)
	emailSplit := strings.Split(req.Email, "@")
	data := map[string]interface{}{
		"name": emailSplit[0],
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

func (svc Service) GetUserByIDs(ids []int) (users []httpAuth.UserResponse, err error) {
	users = []httpAuth.UserResponse{}
	entities, err := svc.repo.GetUserByIDs(ids)
	if err != nil {
		err = errors.Wrap(err, "repo: get user by ids")
		return
	}
	if len(entities) <= 0 {
		return
	}

	for _, entity := range entities {
		var encryptedID string
		encryptedID, err = helper.Encrypt(strconv.Itoa(int(entity.ID)))
		if err != nil {
			err = errors.Wrap(err, "aes: encrypt")
			return
		}
		user := httpAuth.UserResponse{
			ID:    encryptedID,
			Name:  entity.Name,
			Email: entity.Email,
			Role:  string(entity.Role),
		}
		users = append(users, user)
	}

	return
}

func (svc Service) GetUserFromClaims(claims httpAuth.Claims) (resp httpAuth.UserResponse) {
	resp = httpAuth.UserResponse{
		ID:    claims.UserID,
		Name:  claims.UserName,
		Email: claims.UserEmail,
		Role:  claims.UserRole,
	}

	return
}

func (svc Service) handlePagination(entities []dbUser.User, pg helper.Pagination, totalPage int) (newEntities []dbUser.User) {
	if pg.Page == totalPage {
		for i := pg.Offset; i < len(entities); i++ {
			newEntities = append(newEntities, entities[i])
		}
	} else {
		for i := pg.Offset; i < pg.Offset+pg.Limit; i++ {
			newEntities = append(newEntities, entities[i])
		}
	}
	return
}

func (svc Service) GetUsers(pg helper.Pagination) (resps []httpAuth.UserResponse, pgInfo helper.PaginationResponse, err error) {
	entities, err := svc.repo.GetUsers()
	if err != nil {
		err = errors.Wrap(err, "repo: get users")
		return
	}

	totalPage := int(math.Ceil(float64(len(entities)) / float64(pg.Limit)))
	if pg.Page > totalPage {
		err = constant.ErrInvalidPage
		return
	}
	pgInfo = helper.PaginationResponse{
		Page:        pg.Page,
		DataPerPage: pg.Limit,
		TotalPage:   totalPage,
		TotalData:   len(entities),
	}
	newEntities := svc.handlePagination(entities, pg, totalPage)

	for _, entity := range newEntities {
		var encryptedID string
		encryptedID, err = helper.Encrypt(strconv.Itoa(int(entity.ID)))
		if err != nil {
			err = errors.Wrap(err, "aes: encrypt")
			return
		}
		resp := httpAuth.UserResponse{
			ID:    encryptedID,
			Name:  entity.Name,
			Email: entity.Email,
			Role:  string(entity.Role),
		}
		resps = append(resps, resp)
	}

	return
}

func (svc Service) DeleteUser(id int) (err error) {
	// check if user exists
	entities, err := svc.repo.GetUserByIDs([]int{id})
	if err != nil {
		err = errors.Wrap(err, "repo: get user by ids")
		return
	}
	if len(entities) <= 0 {
		err = constant.ErrDataNotFound
		return
	}

	// delete verification
	err = svc.verifSvc.DeleteByEmail(entities[0].Email)
	if err != nil {
		err = errors.Wrap(err, "verif svc: delete by email")
		return
	}

	// delete user
	err = svc.repo.Delete(int(entities[0].ID))
	if err != nil {
		err = errors.Wrap(err, "repo: delete")
	}

	return
}

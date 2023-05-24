package helper

import (
	"regexp"
	"strings"
	"ta/backend/src/constant"
	"unicode"

	"github.com/gin-gonic/gin"
)

type jsonResponseStruct struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func JSONResponse(ctx *gin.Context, code int, msg string, data interface{}) {
	responseStruct := jsonResponseStruct{
		Status: code,
	}

	if msg != "" {
		responseStruct.Message = msg
	}

	if data != nil {
		responseStruct.Data = data
	}

	ctx.JSON(code, responseStruct)
}

type jsonResponseGetAll struct {
	jsonResponseStruct
	PaginationInfo PaginationResponse `json:"pagination_info"`
}

func JSONResponseGetAll(ctx *gin.Context, code int, msg string, paginationInfo PaginationResponse, data interface{}) {
	responseStruct := jsonResponseGetAll{
		jsonResponseStruct: jsonResponseStruct{
			Status: code,
		},
	}

	if msg != "" {
		responseStruct.Message = msg
	}

	if paginationInfo != (PaginationResponse{}) {
		responseStruct.PaginationInfo = paginationInfo
	}

	if data != nil {
		responseStruct.Data = data
	}

	ctx.JSON(code, responseStruct)
}

func SanitizeEmail(email string) (err error) {
	email = strings.ToLower(email)
	if len(email) > constant.MaxEmailLength {
		err = constant.ErrEmailLength
		return
	}

	emailRegex, _ := regexp.Compile(constant.EmailRegex)
	if !emailRegex.MatchString(email) || !isASCII(email) {
		err = constant.ErrEmailInvalid
	}

	return nil
}

func SanitizeUsername(username string) error {
	username = strings.ToLower(username)

	if len(username) < constant.MinUsernameLength || len(username) > constant.MaxUsernameLength {
		return constant.ErrUsernameLength
	}

	usernameRegex, _ := regexp.Compile(constant.UsernameRegex)
	allowedCharacters := usernameRegex.MatchString(username)
	if !allowedCharacters || !isASCII(username) {
		return constant.ErrUsernameUnallowed
	}

	if !unicode.IsLetter(rune(username[0])) {
		return constant.ErrUsernamePrefix
	}

	if strings.HasSuffix(username, "-") ||
		strings.HasSuffix(username, "_") {
		return constant.ErrUsernameSuffix
	}

	return nil
}

func SanitizePassword(password string) error {
	if len(password) < constant.MinPasswordLength || len(password) > constant.MaxPasswordLength {
		return constant.ErrPasswordLength
	}

	hasLetter := false
	hasNumber := false
	hasSpace := false

	for _, i := range password {
		if i > unicode.MaxASCII {
			return constant.ErrPasswordUnallowed
		}

		if unicode.IsLetter(i) {
			hasLetter = true
		} else if unicode.IsDigit(i) {
			hasNumber = true
		} else if unicode.IsSpace(i) {
			hasSpace = true
		}
	}

	if !hasLetter || !hasNumber {
		return constant.ErrPasswordWeak
	}
	if hasSpace {
		return constant.ErrPasswordUnallowed
	}

	return nil
}

func isASCII(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] > unicode.MaxASCII {
			return false
		}
	}
	return true
}

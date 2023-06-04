package helper

import (
	"backend/src/constant"
	httpAuth "backend/src/entity/v1/http/auth"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

func GetJWT(ctx *gin.Context) (token string, err error) {
	authorizationHeader := ctx.GetHeader("Authorization")
	var coookieToken string
	if len(authorizationHeader) > 0 {
		splittedHeader := strings.Split(authorizationHeader, " ")
		if len(splittedHeader) != 2 {
			err = constant.ErrInvalidFormat
			return
		}
		token = splittedHeader[1]
	} else {
		coookieToken, err = ctx.Cookie("Token")
		if err != nil || coookieToken == "" {
			err = constant.ErrTokenNotFound
			return
		}
		token = coookieToken
	}

	return
}

func ExtractJWT(token string) (claims httpAuth.Claims, err error) {
	_ = godotenv.Load()
	claim := &httpAuth.Claims{}
	parsedToken, err := jwt.ParseWithClaims(token, claim, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_ACCESS_SIGNATURE_KEY")), nil
	})
	if err != nil {
		err = jwt.ErrSignatureInvalid
		return
	}
	if !parsedToken.Valid {
		err = constant.ErrTokenInvalid
		return
	}

	claims = *claim

	return
}

func GetClaims(ctx *gin.Context) (claims httpAuth.Claims, err error) {
	token, err := GetJWT(ctx)
	if err != nil {
		return
	}

	claims, err = ExtractJWT(token)
	if err != nil {
		return
	}

	return
}

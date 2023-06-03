package auth

import (
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID       string `json:"id"`
	UserName     string `json:"name"`
	UserEmail    string `json:"email"`
	UserUsername string `json:"username"`
}

type LoginRequest struct {
	Identifier string `json:"identifier" validate:"required"`
	Password   string `json:"password" validate:"required"`
}

type RegisterRequest struct {
	Name     string `json:"name" validate:"required"`
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`

	HashedPassword string `json:"-"`
	PasswordSalt   string `json:"-"`
}

type VerifyRequest struct {
	Username string `form:"username" validate:"required"`
	Token    string `form:"token" validate:"required"`
}

type ResendRequest struct {
	Email string `json:"email" validate:"required"`
}

type DeleteRequest struct {
	ID string `json:"id"`
}

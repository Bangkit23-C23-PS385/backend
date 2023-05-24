package constant

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidFormat          = errors.New("invalid format")
	ErrAccountNotFound        = errors.New("account not found")
	ErrAccountNotVerified     = errors.New("account not verfied")
	ErrAccountAlreadyVerified = errors.New("account already verified")
	ErrInvalidCreds           = errors.New("invalid credentials")

	ErrEmailLength       = fmt.Errorf("email must not exceed %v characters", MaxEmailLength)
	ErrEmailInvalid      = errors.New("invalid email format")
	ErrEmailTaken        = errors.New("email already taken")
	ErrUsernameLength    = fmt.Errorf("must be %d to %d characters long", MinUsernameLength, MaxUsernameLength)
	ErrUsernameUnallowed = errors.New("can only use letters, numbers, underscores, dashes")
	ErrUsernamePrefix    = errors.New("must start with a letter")
	ErrUsernameSuffix    = errors.New("cannot end with dash, or underscore")
	ErrPasswordLength    = fmt.Errorf("must be %d to %d characters long", MinPasswordLength, MaxPasswordLength)
	ErrPasswordUnallowed = errors.New("cannot contain unallowed characters")
	ErrPasswordWeak      = errors.New("must contain letter and number")

	ErrTokenNotFound = errors.New("token not found")
	ErrTokenInvalid  = errors.New("token invalid")

	ErrNoResendAttemptLeft = errors.New("no resend attempt left. please create a new account")

	ErrTemplateNoValue = errors.New("no value detected in template")

	ErrInvalidCipherTextLength = errors.New("invalid ciphertext block size")

	ErrDataNotFound = errors.New("data not found")
	ErrInvalidID    = errors.New("invalid id")
	ErrInvalidPage  = errors.New("invalid page")
)

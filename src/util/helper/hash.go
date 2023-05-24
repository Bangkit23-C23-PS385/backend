package helper

import (
	"crypto/sha256"
	"fmt"

	"github.com/segmentio/ksuid"
)

func HashPassword(password, salt string) string {
	sum := sha256.Sum256([]byte(password + salt))
	return fmt.Sprintf("%x", sum)
}

func GenerateSalt() string {
	salt := ksuid.New()
	return salt.String()
}

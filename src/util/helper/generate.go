package helper

import (
	"math/rand"
	"time"

	"github.com/rs/xid"
)

func GenerateEmailVerifToken() (token string) {
	token = xid.New().String()
	return
}

// Generate 16-letters random string
func GenereateRandomString() (str string) {
	rand.Seed(time.Now().UnixNano())
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	pass := make([]rune, 16)
	for i := range pass {
		pass[i] = letters[rand.Intn(len(letters))]
	}
	str = string(pass)
	return
}

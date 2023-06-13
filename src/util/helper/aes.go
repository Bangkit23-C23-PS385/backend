package helper

import (
	"backend/src/constant"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
	"os"

	"github.com/joho/godotenv"
	"github.com/pkg/errors"
)

func Encrypt(message string) (encryptedMessage string, err error) {
	_ = godotenv.Load()
	key := []byte(os.Getenv("AES_KEY"))
	plainText := []byte(message)
	block, err := aes.NewCipher(key)
	if err != nil {
		err = errors.Wrap(err, "aes encrypt: cipher block creation")
		return
	}

	cipherText := make([]byte, aes.BlockSize+len(plainText))
	iv := cipherText[:aes.BlockSize]
	_, err = io.ReadFull(rand.Reader, iv)
	if err != nil {
		err = errors.Wrap(err, "aes encrypt: rand generation")
		return
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], plainText)
	encryptedMessage = base64.URLEncoding.EncodeToString(cipherText)

	return
}

func Decrypt(encryptedMessage string) (message string, err error) {
	_ = godotenv.Load()
	key := []byte(os.Getenv("AES_KEY"))
	cipherText, err := base64.URLEncoding.DecodeString(encryptedMessage)
	if err != nil {
		err = errors.Wrap(err, "aes decrypt: could not decode message")
		return
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		err = errors.Wrap(err, "aes encrypt: cipher block creation")
		return
	}

	if len(cipherText) < aes.BlockSize {
		err = constant.ErrInvalidCipherTextLength
		return
	}

	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(cipherText, cipherText)

	message = string(cipherText)

	return
}

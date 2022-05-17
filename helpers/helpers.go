package helpers

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/imayrus/jwt-api/models"
	"golang.org/x/crypto/bcrypt"
)

var secretkey string = "jwtsecretkey"

func SetError(err models.Error, message string) models.Error {
	err.IsError = true
	err.Message = message
	return err
}

func GenerateHashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateJWT(email, role, firstname string) (string, error) {
	var SigningKey = []byte(secretkey)
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["authorized"] = true
	claims["email"] = email
	claims["firstname"] = firstname
	// claims["lastname"] = lastname
	claims["role"] = role
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	TokenString, err := token.SignedString(SigningKey)
	if err != nil {
		fmt.Println("Something went wrong")
		return "", err
	}

	return TokenString, nil
}

package accounts

import (
	"golang.org/x/crypto/bcrypt"
)

type Account struct {
	Username     string `json:"username"`
	PasswordHash string `json:"password_hash"`
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

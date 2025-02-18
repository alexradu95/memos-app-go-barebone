package auth

import (
	"database/sql"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const secret = "your-secret-key" //  Use environment variable in production

type Account struct {
	Id           string
	Username     string
	PasswordHash string
}

type MyCustomClaims struct {
	UserID string `json:"userId"`
	jwt.RegisteredClaims
}

func Login(db *sql.DB, username string, password string) (string, error) {
	var account Account

	err := db.QueryRow("SELECT id, username, password_hash FROM accounts WHERE username = ?", username).
		Scan(&account.Id, &account.Username, &account.PasswordHash)
	if err != nil {
		return "", errors.New("Invalid username or password.")
	}

	err = bcrypt.CompareHashAndPassword([]byte(account.PasswordHash), []byte(password))
	if err != nil {
		return "", errors.New("Invalid username or password.")
	}

	token, err := generateToken(account)
	if err != nil {
		return "", errors.New("Error occurred while generating the token.")
	}

	return token, nil
}

func generateToken(account Account) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour) // Token valid for 24 hours

	claims := MyCustomClaims{
		UserID: account.Id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "journal-lite",
			Subject:   account.Username, //  Use username as subject
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secret)) //  byte slice of secret
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func ValidateToken(tokenString string) (*MyCustomClaims, error) {
	claims := &MyCustomClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil //  Return the secret key
	})

	if err != nil {
		return nil, err // Return specific error
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

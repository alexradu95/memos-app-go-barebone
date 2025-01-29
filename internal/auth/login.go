package auth

import (
	"database/sql"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

const secret = "your-secret-key"

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	IsSuccess    bool   `json:"isSuccess"`
	BearerToken  string `json:"bearerToken"`
	RefreshToken string `json:"refreshToken"`
}

type User struct {
	Id       string `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func LoginHandler(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req LoginRequest
		var user User

		err := c.Bind(&req)
		if err != nil {
			return c.JSON(400, map[string]string{
				"message": "Invalid request.",
			})
		}

		err = db.QueryRow("SELECT id, email, password FROM users WHERE email = $1", req.Email).
			Scan(&user.Id, &user.Email, &user.Password)

		if err != nil && err == sql.ErrNoRows {
			return c.JSON(401, map[string]string{
				"message": "Invalid email or password.",
			})
		}

		if err != nil {
			return c.JSON(500, map[string]string{
				"message": "Error occurred while querying the database.",
			})
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
		if err != nil {
			return c.JSON(401, map[string]string{
				"message": "Invalid email or password.",
			})
		}

		bearerToken, err := generateBearerToken(user)
		if err != nil {
			return c.JSON(500, map[string]string{
				"message": "Error occurred while generating the bearer token.",
			})
		}

		refreshToken, err := generateRefreshToken(user)
		if err != nil {
			return c.JSON(500, map[string]string{
				"message": "Error occurred while generating the refresh token.",
			})
		}

		res := LoginResponse{
			IsSuccess:    true,
			BearerToken:  bearerToken,
			RefreshToken: refreshToken,
		}

		return c.JSON(200, res)
	}
}

func generateBearerToken(user User) (string, error) {
	claims := jwt.MapClaims{
		"userId": user.Id,
		"email":  user.Email,
		"exp":    time.Now().Add(time.Minute * 5).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	bearerToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return bearerToken, nil
}

func generateRefreshToken(user User) (string, error) {
	claims := jwt.MapClaims{
		"userId": user.Id,
		"email":  user.Email,
		"exp":    time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	refreshToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return refreshToken, nil
}

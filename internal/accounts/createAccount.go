package accounts

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type CreateAccountRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateAccountResponse struct {
	IsSuccess bool   `json:"isSuccess"`
	Message   string `json:"message"`
}

func CreateAccountHandler(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req CreateAccountRequest

		err := c.Bind(&req)
		if err != nil {
			res := CreateAccountResponse{
				IsSuccess: false,
				Message:   "Invalid request.",
			}
			return c.JSON(400, res)
		}

		var count int
		err = db.QueryRow("SELECT COUNT(*) FROM users WHERE email = $1", req.Email).
			Scan(&count)
		if err != nil {
			res := CreateAccountResponse{
				IsSuccess: false,
				Message:   "Error occured while checking if account exists.",
			}
			return c.JSON(500, res)
		}

		if count != 0 {
			response := CreateAccountResponse{
				IsSuccess: false,
				Message:   "Account with that email already exists.",
			}
			return c.JSON(409, response)
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			response := CreateAccountResponse{
				IsSuccess: false,
				Message:   "Error occured while hashing password.",
			}
			return c.JSON(500, response)
		}

		_, err = db.Exec(
			"INSERT INTO users (id, email, password) VALUES ($1, $2, $3)",
			uuid.New().String(),
			req.Email,
			hashedPassword,
		)
		if err != nil {
			res := CreateAccountResponse{
				IsSuccess: false,
				Message:   "Error occured while creating account.",
			}
			return c.JSON(500, res)
		}

		res := CreateAccountResponse{
			IsSuccess: true,
			Message:   "Account created successfully.",
		}

		return c.JSON(201, res)
	}
}

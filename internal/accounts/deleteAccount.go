package accounts

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type DeleteAccountResponse struct {
	IsSuccess bool   `json:"isSuccess"`
	Message   string `json:"message"`
}

func DeleteAccountHandler(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {

		_, err := db.Exec(
			"DELETE FROM users WHERE email = $1",
			uuid.New().String(),
		)
		if err != nil {
			res := DeleteAccountResponse{
				IsSuccess: false,
				Message:   "Error occured while deleting account.",
			}
			return c.JSON(500, res)
		}

		res := DeleteAccountResponse{
			IsSuccess: true,
			Message:   "Account deleted successfully.",
		}

		return c.JSON(201, res)
	}
}

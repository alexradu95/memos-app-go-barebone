package accounts

import (
	"database/sql"
	"journal-lite/internal/database"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Account struct {
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`
}

func CreateAccountHandler(db *sql.DB, newAccount Account) (int64, error) {

	count, err := RetrieveCountOfAccountsWithEmail(newAccount.Email)

	if count != 0 {
		return 0, nil
	}

	hashedPassword, err := HashPassword(newAccount.PasswordHash)

	if err != nil {
		return 0, err
	}

	newAccount.PasswordHash = hashedPassword

	err = AddAccountToDatabase(newAccount)

	if err != nil {
		return 0, err
	}

	return 1, nil
}

func RetrieveCountOfAccountsWithEmail(email string) (int, error) {
	var count int
	err := database.Db.QueryRow("SELECT COUNT(*) FROM accounts WHERE email = $1", email).
		Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func AddAccountToDatabase(newAccount Account) error {
	_, err := database.Db.Exec(
		"INSERT INTO users (id, email, password) VALUES ($1, $2, $3)",
		uuid.New().String(),
		newAccount.Email,
		newAccount.PasswordHash,
	)
	if err != nil {
		return err
	}
	return nil
}

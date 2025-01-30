package accounts

import (
	"database/sql"
	"journal-lite/internal/database"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Account struct {
	Username     string `json:"username"`
	PasswordHash string `json:"password_hash"`
}

func CreateAccount(db *sql.DB, newAccount Account) (int64, error) {

	count, err := RetrieveCountOfAccountsWithUsername(newAccount.Username)

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

func RetrieveCountOfAccountsWithUsername(username string) (int, error) {
	var count int
	err := database.Db.QueryRow("SELECT COUNT(*) FROM accounts WHERE username = ?", username).
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
		"INSERT INTO accounts (username, password_hash, created_at) VALUES (?, ?, ?)",
		newAccount.Username,
		newAccount.PasswordHash,
		time.Now(),
	)
	if err != nil {
		return err
	}
	return nil
}

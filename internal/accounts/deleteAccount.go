package accounts

import (
	"database/sql"
)

func DeleteAccountById(db *sql.DB, accountId int64) error {
	_, err := db.Exec(
		"DELETE FROM accounts WHERE id = $1",
		accountId,
	)
	if err != nil {
		return err
	}

	return nil
}

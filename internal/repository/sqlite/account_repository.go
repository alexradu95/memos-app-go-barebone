// internal/repository/sqlite/account_repository.go
package sqlite

import (
	"context"
	"database/sql"
	"journal-lite/internal/accounts"
	"journal-lite/internal/repository"
	"time"
)

type AccountRepository struct {
	db *sql.DB
}

func NewAccountRepository(db *sql.DB) repository.AccountRepository {
	return &AccountRepository{db: db}
}

func (r *AccountRepository) CreateAccount(ctx context.Context, account accounts.Account) (int64, error) {
	count, err := r.RetrieveCountOfAccountsWithUsername(ctx, account.Username)
	if err != nil {
		return 0, err
	}
	if count != 0 {
		return 0, nil
	}

	hashedPassword, err := accounts.HashPassword(account.PasswordHash)
	if err != nil {
		return 0, err
	}

	account.PasswordHash = hashedPassword

	_, err = r.db.ExecContext(ctx,
		"INSERT INTO accounts (username, password_hash, created_at) VALUES (?, ?, ?)",
		account.Username,
		account.PasswordHash,
		time.Now(),
	)
	if err != nil {
		return 0, err
	}

	return 1, nil
}

func (r *AccountRepository) DeleteAccountById(ctx context.Context, accountId int64) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM accounts WHERE id = ?", accountId)
	return err
}

func (r *AccountRepository) RetrieveCountOfAccountsWithUsername(ctx context.Context, username string) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM accounts WHERE username = ?", username).Scan(&count)
	return count, err
}

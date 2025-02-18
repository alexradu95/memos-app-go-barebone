// internal/repository/account_repository.go
package repository

import (
	"context"
	"journal-lite/internal/accounts"
)

type AccountRepository interface {
	CreateAccount(ctx context.Context, account accounts.Account) (int64, error)
	DeleteAccountById(ctx context.Context, accountId int64) error
	RetrieveCountOfAccountsWithUsername(ctx context.Context, username string) (int, error)
}

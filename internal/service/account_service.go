// internal/service/account_service.go
package service

import (
	"context"
	"journal-lite/internal/accounts"
	"journal-lite/internal/repository"
)

type AccountService struct {
	repo repository.AccountRepository
}

func NewAccountService(repo repository.AccountRepository) *AccountService {
	return &AccountService{repo: repo}
}

func (s *AccountService) CreateAccount(ctx context.Context, account accounts.Account) (int64, error) {
	return s.repo.CreateAccount(ctx, account)
}

func (s *AccountService) DeleteAccountById(ctx context.Context, accountId int64) error {
	return s.repo.DeleteAccountById(ctx, accountId)
}

// internal/service/post_service.go

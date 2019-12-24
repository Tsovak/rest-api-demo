package service

import (
	"context"
	pg "github.com/go-pg/pg/v9"
	"github.com/pkg/errors"
	"github.com/tsovak/rest-api-demo/api/model"
	"github.com/tsovak/rest-api-demo/repositories"
)

//go:generate mockgen -source account_manager.go -package mock -destination ../mock/account_manager.go

// AccountManager declare interface to access an accounts
type AccountManager interface {
	// GetAllAccounts return all accounts
	GetAllAccounts(ctx context.Context) ([]model.Account, error)

	// Save new account
	Save(ctx context.Context, account *model.Account) error

	// Find account by ID
	FindByID(ctx context.Context, id string) (model.Account, error)

	// Delete account by ID
	DeleteByID(ctx context.Context, id string) error

	// Update account
	Update(ctx context.Context, account *model.Account, fn func(*pg.Tx) error) error
}

type accountManager struct {
	accountRepository repositories.AccountRepository
}

// NewAccountManager returns interface to access an account
func NewAccountManager(accountRepository repositories.AccountRepository) AccountManager {
	return &accountManager{
		accountRepository: accountRepository,
	}
}

func (m *accountManager) GetAllAccounts(ctx context.Context) ([]model.Account, error) {
	accounts, err := m.accountRepository.GetAll(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get all accounts")
	}
	return accounts, nil
}

func (m *accountManager) Save(ctx context.Context, account *model.Account) error {
	err := m.accountRepository.Save(ctx, account)
	if err != nil {
		return errors.Wrap(err, "Failed to save account")
	}
	return nil
}

func (m *accountManager) FindByID(ctx context.Context, id string) (model.Account, error) {
	accounts, err := m.accountRepository.FindByID(ctx, id)
	if err != nil {
		return model.Account{}, errors.Wrap(err, "Failed to find account")
	}
	return accounts, nil
}

func (m *accountManager) DeleteByID(ctx context.Context, id string) error {
	err := m.accountRepository.DeleteByID(ctx, id)
	if err != nil {
		return errors.Wrap(err, "Failed to delete account")
	}
	return nil
}

func (m *accountManager) Update(ctx context.Context, account *model.Account, fn func(*pg.Tx) error) error {
	err := m.accountRepository.Update(ctx, account, fn)
	if err != nil {
		return errors.Wrap(err, "Failed to update account")
	}
	return nil
}

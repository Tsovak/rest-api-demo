package service

import (
	"context"
	"github.com/pkg/errors"
	"github.com/tsovak/rest-api-demo/api/model"
	"github.com/tsovak/rest-api-demo/repositories"
)

//go:generate mockgen -source account_manager.go -package mock -destination ../mock/account_manager.go

// AccountManager declare interface to access accounts
type AccountManager interface {
	// GetAllAccounts return all accounts
	GetAllAccounts(ctx context.Context) ([]model.Account, error)

	// Store save new account
	Save(ctx context.Context, account *model.Account) error

	// Find account by id
	FindById(ctx context.Context, id string) (model.Account, error)

	// Delete account by id
	DeleteById(ctx context.Context, id string) error

	// Update account
	Update(ctx context.Context, account *model.Account) error
}

type accountManager struct {
	accountRepository repositories.AccountRepository
}

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

func (m *accountManager) FindById(ctx context.Context, id string) (model.Account, error) {
	accounts, err := m.accountRepository.FindById(ctx, id)
	if err != nil {
		return model.Account{}, errors.Wrap(err, "Failed to find account")
	}
	return accounts, nil
}

func (m *accountManager) DeleteById(ctx context.Context, id string) error {
	err := m.accountRepository.DeleteById(ctx, id)
	if err != nil {
		return errors.Wrap(err, "Failed to delete account")
	}
	return nil
}

func (m *accountManager) Update(ctx context.Context, account *model.Account) error {
	err := m.accountRepository.Update(ctx, account)
	if err != nil {
		return errors.Wrap(err, "Failed to update account")
	}
	return nil
}

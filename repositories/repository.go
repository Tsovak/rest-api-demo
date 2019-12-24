package repositories

import (
	"context"
	pg "github.com/go-pg/pg/v9"
	"github.com/tsovak/rest-api-demo/api/model"
)

//go:generate mockgen -source repository.go -package mock -destination ../mock/repository.go

// AccountRepository declare repository for accounts
type AccountRepository interface {
	// GetAll return all accounts in storage
	GetAll(ctx context.Context) ([]model.Account, error)

	// Save new account in storage
	Save(ctx context.Context, account *model.Account) error

	// Find Account by ID
	FindByID(ctx context.Context, id string) (model.Account, error)

	// Delete Account by ID
	DeleteByID(ctx context.Context, id string) error

	// Update account
	Update(ctx context.Context, account *model.Account, fn func(tx *pg.Tx) error) error
}

// PaymentRepository declare repository for payments
type PaymentRepository interface {
	// GetAll return all payments
	GetAll(ctx context.Context) ([]model.Payment, error)

	// Save new payment
	Save(ctx context.Context, payment ...*model.Payment) error

	// Find Payment by ID
	FindByID(ctx context.Context, id string) (model.Payment, error)

	// Get Payments by account ID
	GetPaymentsByAccountID(ctx context.Context, accountID string) ([]model.Payment, error)

	// Get function for save payments without commit in DB
	GetSaveTransaction(ctx context.Context, payments ...*model.Payment) func(tx *pg.Tx) error
}

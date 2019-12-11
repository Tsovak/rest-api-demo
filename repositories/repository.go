package repositories

import (
	"context"
	"github.com/tsovak/rest-api-demo/api/model"
)

//go:generate mockgen -source repository.go -package mock -destination ../mock/repository.go

// AccountRepository declare repository for accounts
type AccountRepository interface {
	// GetAll return all accounts in storage
	GetAll(ctx context.Context) ([]model.Account, error)

	// Save new account in storage
	Save(ctx context.Context, account *model.Account) error

	// Find Account by id
	FindById(ctx context.Context, id string) (model.Account, error)

	// Delete Account by id
	DeleteById(ctx context.Context, id string) error

	// Update account
	Update(ctx context.Context, account *model.Account) error
}

// PaymentRepository declare repository for payments
type PaymentRepository interface {
	// GetAll return all payments
	GetAll(ctx context.Context) ([]model.Payment, error)

	// Save new payment
	Save(ctx context.Context, payment *model.Payment) error

	// Find Payment by id
	FindById(ctx context.Context, id string) (model.Payment, error)
}

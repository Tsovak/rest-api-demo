package service

import (
	"context"
	pg "github.com/go-pg/pg/v9"
	"github.com/pkg/errors"
	"github.com/tsovak/rest-api-demo/api/model"
	"github.com/tsovak/rest-api-demo/repositories"
)

//go:generate mockgen -source payment_manager.go -package mock -destination ../mock/payment_manager.go

// PaymentManager declare interface to access accounts
type PaymentManager interface {
	// GetAccountPayments return all payments
	GetAccountPayments(ctx context.Context) ([]model.Payment, error)

	// CreatePayments save the payments
	CreatePayments(ctx context.Context, payments ...*model.Payment) error

	// GetPaymentsByAccountId return all payments fot specified account
	GetPaymentsByAccountId(ctx context.Context, id string) ([]model.Payment, error)

	// GetSaveTransaction return transaction function for save in DB without commit
	GetSaveTransaction(ctx context.Context, payments ...*model.Payment) func(*pg.Tx) error
}

type paymentManager struct {
	paymentRepository repositories.PaymentRepository
}

func NewPaymentManager(repository repositories.PaymentRepository) PaymentManager {
	return &paymentManager{
		paymentRepository: repository,
	}
}

func (p *paymentManager) GetAccountPayments(ctx context.Context) ([]model.Payment, error) {
	payments, err := p.paymentRepository.GetAll(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get all payments ")
	}
	return payments, nil
}

func (p *paymentManager) CreatePayments(ctx context.Context, payments ...*model.Payment) error {
	err := p.paymentRepository.Save(ctx, payments...)
	if err != nil {
		return errors.Wrap(err, "Failed to save payment")
	}
	return nil
}

func (p *paymentManager) GetPaymentsByAccountId(ctx context.Context, accountId string) ([]model.Payment, error) {
	if len(accountId) == 0 {
		return nil, errors.New("Account id is incorrect")
	}

	payments, err := p.paymentRepository.GetPaymentsByAccountId(ctx, accountId)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get payments for an account")
	}

	return payments, nil
}

func (p *paymentManager) GetSaveTransaction(ctx context.Context, payments ...*model.Payment) func(tx *pg.Tx) error {
	return p.paymentRepository.GetSaveTransaction(ctx, payments...)
}

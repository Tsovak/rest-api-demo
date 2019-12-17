package service

import (
	"context"
	"github.com/pkg/errors"
	"github.com/tsovak/rest-api-demo/api/model"
	"github.com/tsovak/rest-api-demo/repositories"
)

//go:generate mockgen -source payment_manager.go -package mock -destination ../mock/payment_manager.go

// PaymentManager declare interface to access accounts
type PaymentManager interface {
	// GetAllPayments return all payments
	GetAllPayments(ctx context.Context) ([]model.Payment, error)
	CreatePayment(ctx context.Context, payment *model.Payment) error
}

type paymentManager struct {
	paymentRepository repositories.PaymentRepository
}

func NewPaymentManager(repository repositories.PaymentRepository) PaymentManager {
	return &paymentManager{
		paymentRepository: repository,
	}
}

func (p *paymentManager) GetAllPayments(ctx context.Context) ([]model.Payment, error) {
	payments, err := p.paymentRepository.GetAll(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get all payments ")
	}
	return payments, nil
}

func (p *paymentManager) CreatePayment(ctx context.Context, payment *model.Payment) error {
	err := p.paymentRepository.Save(ctx, payment)
	if err != nil {
		return errors.Wrap(err, "Failed to save payment")
	}
	return nil
}

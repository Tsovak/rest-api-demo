package repositories

import (
	"context"
	"github.com/go-pg/pg/v9"
	"github.com/pkg/errors"
	"github.com/tsovak/rest-api-demo/api/model"
	"strconv"
)

type paymentRepository struct {
	Db *pg.DB
}

func NewPaymentRepository(db *pg.DB) PaymentRepository {
	return &paymentRepository{
		Db: db,
	}
}

func (p paymentRepository) GetAll(ctx context.Context) ([]model.Payment, error) {
	var payments []model.Payment
	err := p.Db.WithContext(ctx).Model(&payments).Select()
	if err != nil {
		return nil, err
	}

	if payments == nil {
		return []model.Payment{}, nil
	}
	return payments, nil
}

func (p paymentRepository) Save(ctx context.Context, payment *model.Payment) error {
	if payment == nil {
		return errors.New("Input parameter payment is nil")
	}

	result, err := p.Db.WithContext(ctx).Model(payment).Returning("id").Insert(&payment.ID)
	if err != nil {
		return errors.Wrapf(err, "Failed to insert payment %v", payment)
	}

	if result != nil {
		if result.RowsAffected() == 0 {
			return errors.New("Failed to insert, affected is 0")
		}
	}

	return nil
}

func (p paymentRepository) FindById(ctx context.Context, id string) (model.Payment, error) {
	if len(id) == 0 {
		return model.Payment{}, errors.New("Payment id is incorrect")
	}

	index, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return model.Payment{}, errors.Wrap(err, "Cannot convert Payment id")
	}
	payment := model.Payment{ID: index}
	err = p.Db.WithContext(ctx).Select(&payment)
	//pg returns error when no rows in the result set
	if err == pg.ErrNoRows {
		// if row is empty than return empty model
		return model.Payment{}, nil
	}
	return payment, err
}

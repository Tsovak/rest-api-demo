package repositories

import (
	"context"
	"fmt"
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

func (p paymentRepository) GetSaveTransaction(ctx context.Context, payments ...*model.Payment) func(*pg.Tx) error {
	fn := func(tx *pg.Tx) error {
		if payments == nil {
			return errors.New("Input parameter payment is nil")
		}

		if tx == nil {
			return errors.New("Transaction is nil")
		}

		for _, payment := range payments {
			result, err := tx.Model(payment).Returning("id").Insert(&payment.ID)
			if err != nil {
				return errors.Wrapf(err, "Failed to insert payment %v", payment)
			}
			if result != nil {
				if result.RowsAffected() == 0 {
					return errors.New("Failed to insert, affected is 0")
				}
			}
		}
		return nil
	}

	return fn
}

func (p paymentRepository) Save(ctx context.Context, payments ...*model.Payment) error {
	if payments == nil {
		return errors.New("Input parameter payment is nil")
	}

	// we need to save all payments at the same time
	tx, trErr := p.Db.WithContext(ctx).Begin()
	if trErr != nil {
		return errors.Wrapf(trErr, "Cannot start transaction")
	}

	for _, payment := range payments {
		result, err := tx.Model(payment).Returning("id").Insert(&payment.ID)
		if err != nil {
			_ = tx.Rollback()
			return errors.Wrapf(err, "Failed to insert payment %v", payment)
		}
		if result != nil {
			if result.RowsAffected() == 0 {
				_ = tx.Rollback()
				return errors.New("Failed to insert, affected is 0")
			}
		}
	}

	trErr = tx.Commit()
	if trErr != nil {
		return errors.Wrapf(trErr, "Cannot commit transaction")
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

func (p paymentRepository) GetPaymentsByAccountId(ctx context.Context, accountId string) ([]model.Payment, error) {
	if len(accountId) == 0 {
		return nil, errors.New("accountId id is incorrect")
	}

	_, err := strconv.ParseInt(accountId, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "Cannot convert accountId id")
	}

	var localPayments []model.Payment
	err = p.Db.WithContext(ctx).
		Model(&model.Payment{}).
		Where(fmt.Sprintf("to_account_id='%v' or from_account_id='%v'", accountId, accountId)).
		Select(&localPayments)

	// if account has not any payment return empty
	if err == nil && localPayments == nil {
		localPayments = []model.Payment{}
	}

	return localPayments, err
}

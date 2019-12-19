package repositories

import (
	"context"
	"fmt"
	"github.com/go-pg/pg/v9"
	"github.com/pkg/errors"
	"github.com/tsovak/rest-api-demo/api/model"
	"strconv"
)

type accountRepository struct {
	Db *pg.DB
}

func NewAccountRepository(Db *pg.DB) AccountRepository {
	return &accountRepository{
		Db: Db,
	}
}

func (a accountRepository) GetAll(ctx context.Context) ([]model.Account, error) {
	var users []model.Account
	err := a.Db.WithContext(ctx).Model(&users).Select()
	if err != nil {
		return nil, err
	}

	if users == nil {
		return []model.Account{}, nil
	}
	return users, nil
}

func (a accountRepository) Save(ctx context.Context, account *model.Account) error {
	if account == nil {
		return errors.New("Input parameter account is nil")
	}

	result, err := a.Db.WithContext(ctx).Model(account).Returning("id").Insert(&account.ID)
	if err != nil {
		return errors.Wrapf(err, "Failed to insert account %v", account)
	}

	if result != nil {
		if result.RowsAffected() == 0 {
			return errors.New("Failed to insert, affected is 0")
		}
	}

	return nil
}

func (a accountRepository) FindById(ctx context.Context, id string) (model.Account, error) {
	if len(id) == 0 {
		return model.Account{}, errors.New("Account id is incorrect")
	}

	index, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return model.Account{}, errors.Wrap(err, "Cannot convert Account id")
	}
	account := model.Account{ID: index}
	err = a.Db.WithContext(ctx).Select(&account)
	//pg returns error when no rows in the result set
	if err == pg.ErrNoRows {
		// if row is empty than return empty model
		return model.Account{}, nil
	}
	return account, err
}

func (a accountRepository) DeleteById(ctx context.Context, id string) error {
	if len(id) == 0 {
		return errors.New("Account id is incorrect")
	}

	index, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return errors.Wrap(err, "Cannot convert Account id")
	}

	account := &model.Account{ID: index}
	return a.Db.WithContext(ctx).Delete(account)
}

func (a accountRepository) Update(ctx context.Context, account *model.Account, fn func(tx *pg.Tx) error) error {
	if account == nil {
		return errors.New("Account is nil")
	}

	// prepare the update function which will be executed later in one tx
	updateTransactionFn := func(tx *pg.Tx) error {
		err := tx.Update(account)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("Cannot update balance for account with id %v", account.ID))
		}
		return nil
	}

	// start a transaction
	tx, err := a.Db.WithContext(ctx).Begin()
	if err != nil {
		return errors.Wrap(err, "Cannot start transaction")
	}

	// if additional update function is not nil we need to execute it
	if fn != nil {
		// do update
		err = fn(tx)
		if err != nil {
			_ = tx.Rollback()
			_ = tx.Close()
			return err
		}
	}

	// do update account
	err = updateTransactionFn(tx)
	if err != nil {
		_ = tx.Rollback()
		_ = tx.Close()
		return err
	}

	err = tx.Commit()
	if err != nil {
		_ = tx.Close()
		return err
	}

	return nil
}

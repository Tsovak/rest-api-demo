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
	return a.Db.WithContext(ctx).Insert(account)
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

func (a accountRepository) Update(ctx context.Context, account *model.Account) error {
	if account == nil {
		return errors.New("Account is nil")
	}

	err := a.Db.WithContext(ctx).Update(account)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Cannot update balance for account with id %v", account.ID))
	}
	return nil
}

package repositories

import (
	"context"
	"github.com/go-pg/pg/v9"
	"github.com/tsovak/rest-api-demo/api/model"
)

type accountRepository struct {
	Db *pg.DB
}

func NewAccountRepository(Db *pg.DB) AccountRepository {
	return &accountRepository{Db: Db}
}

func (a accountRepository) GetAll(ctx context.Context) ([]model.Account, error) {
	panic("implement me")
}

func (a accountRepository) Save(ctx context.Context, account *model.Account) error {
	panic("implement me")
}

func (a accountRepository) FindById(ctx context.Context, id string) (model.Account, error) {
	panic("implement me")
}

func (a accountRepository) DeleteById(ctx context.Context, id string) (model.Account, error) {
	panic("implement me")
}

func (a accountRepository) UpdateBalance(ctx context.Context, id string, incr int64) error {
	panic("implement me")
}

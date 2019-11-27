package testutils

import (
	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
	"github.com/tsovak/rest-api-demo/api/model"
	"math/rand"
)

func GetTestUser() *model.Account {
	user := &model.Account{
		ID: func() int64 {
			return rand.Int63()
		}(),
		Name:     "Test name",
		Currency: "RUB",
		Balance:  100,
	}

	return user
}

func CreateSchema(db *pg.DB) error {
	for _, model := range []interface{}{&model.Payment{}, &model.Account{}} {
		err := db.CreateTable(model, &orm.CreateTableOptions{
			Temp:          true,
			FKConstraints: true,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

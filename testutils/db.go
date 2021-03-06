package testutils

import (
	"fmt"
	"github.com/go-pg/migrations/v7"
	"github.com/go-pg/pg/v9"
	"github.com/ory/dockertest/v3"
	"github.com/pkg/errors"
	"github.com/tsovak/rest-api-demo/api/model"
	"math/rand"
)

// GetTestUser returns the new random test User
func GetTestUser() *model.Account {
	user := &model.Account{
		ID:       rand.Int63(),
		Name:     "Test name",
		Currency: "RUB",
		Balance:  100,
	}

	return user
}

// GetTestPayment returns a new random payment
func GetTestPayment() *model.Payment {
	return &model.Payment{
		ID:            rand.Int63(),
		Amount:        100,
		ToAccountID:   "123456",
		FromAccountID: "654321",
		Direction:     model.Outgoing,
	}
}

type DBSetup struct {
	Db        *pg.DB
	PgOptions pg.Options
	Cleaner   func() error
}

// SetupTestDB performs a db and migration for integration tests
func SetupTestDB(pgOptions *pg.Options, migrationsDir string) (*DBSetup, error) {
	var err error

	pool, err := dockertest.NewPool("")
	if err != nil {
		return nil, errors.Wrap(err, "Could not connect to docker")
	}

	resource, err := pool.Run(
		"postgres", "12",
		[]string{
			"POSTGRES_DB=" + pgOptions.Database,
			"POSTGRES_PASSWORD=" + pgOptions.Password,
		},
	)
	if err != nil {
		return nil, errors.Wrap(err, "Could not start resource")
	}

	poolCleaner := func() error {
		// When you're done, kill and remove the container
		err := pool.Purge(resource)
		if err != nil {
			return errors.Wrap(err, "failed to purge docker pool")
		}
		return nil
	}

	options := *pgOptions
	options.Addr = fmt.Sprintf("%s:%s", options.Addr, resource.GetPort("5432/tcp"))

	var db *pg.DB
	err = pool.Retry(func() error {
		db = pg.Connect(&options)
		_, err := db.Exec("select 1")
		return err
	})
	if err != nil {
		returnedError := errors.Wrap(err, "Could not start postgres")
		_ = poolCleaner()
		errwrap := errors.Wrap(returnedError, "Clean poll error")
		return nil, errors.Wrap(errwrap, "Could not clean db")
	}

	dbCleaner := func() error {
		err := db.Close()
		if err != nil {
			return errors.Wrap(err, "failed to purge docker pool")
		}
		return nil
	}

	cleaner := func() error {
		err := dbCleaner()
		if err != nil {
			return err
		}
		err = poolCleaner()
		if err != nil {
			return err
		}
		return nil
	}

	migrationCollection := migrations.NewCollection()

	_, _, err = migrationCollection.Run(db, "init")
	if err != nil {
		_ = cleaner()
		return nil, errors.Wrap(err, "Could not init migrations")

	}

	err = migrationCollection.DiscoverSQLMigrations(migrationsDir)
	if err != nil {
		_ = cleaner()
		return nil, errors.Wrap(err, "Failed to read migrations")
	}

	_, _, err = migrationCollection.Run(db, "up")
	if err != nil {
		_ = cleaner()
		return nil, errors.Wrap(err, "Could not migrate")
	}
	return &DBSetup{
		Db:        db,
		PgOptions: options,
		Cleaner:   cleaner,
	}, nil
}

package db

import (
	"github.com/sirupsen/logrus"
	"github.com/tsovak/rest-api-demo/config"

	"github.com/go-pg/pg/v9"
)

const (
	ReadTimeout  = 30_000
	WriteTimeout = 30_000
	PoolSize     = 10
	MinIdleConns = 10
)

type postgresClient struct {
	Logger *logrus.FieldLogger
	Db     *pg.DB
}

func (p postgresClient) GetConnection() *pg.DB {
	return p.Db
}

func (p postgresClient) Close() error {
	return p.Db.Close()
}

func NewPostgresClient(logger *logrus.FieldLogger, config config.Config) PostgresClient {
	db := pg.Connect(getPgConnectionOptions(config))
	return postgresClient{
		Logger: logger,
		Db:     db,
	}
}

func getPgConnectionOptions(config config.Config) *pg.Options {
	return &pg.Options{
		Addr:            config.DbConfig.Address,
		User:            config.DbConfig.Username,
		Password:        config.DbConfig.Password,
		Database:        config.DbConfig.Database,
		ApplicationName: "",
		ReadTimeout:     ReadTimeout,
		WriteTimeout:    WriteTimeout,
		PoolSize:        PoolSize,
		MinIdleConns:    MinIdleConns,
	}
}

package db

import (
	"github.com/sirupsen/logrus"
	"github.com/tsovak/rest-api-demo/config"
	"time"

	"github.com/go-pg/pg/v9"
)

const (
	ReadTimeout  = 30 * time.Second
	WriteTimeout = 30 * time.Second
	PoolSize     = 10
	MinIdleConns = 10
)

type postgresClient struct {
	Db *pg.DB
}

func (p postgresClient) GetConnection() *pg.DB {
	return p.Db
}

func (p postgresClient) Close() error {
	return p.Db.Close()
}

func NewPostgresClientFromConfig(config config.Config) PostgresClient {
	return NewPostgresClientFromPgOptions(config.Logger, GetPgConnectionOptions(config))
}

func NewPostgresClientFromPgOptions(logger *logrus.Logger, pgOptions *pg.Options) PostgresClient {
	logrus.Debug("Trying to connect " + pgOptions.Addr)
	db := pg.Connect(pgOptions)
	return postgresClient{
		Db: db,
	}
}
func NewPostgresClient(db *pg.DB) PostgresClient {
	return postgresClient{
		Db: db,
	}
}

func GetPgConnectionOptions(config config.Config) *pg.Options {
	return &pg.Options{
		Addr:            config.DbConfig.Address,
		User:            config.DbConfig.Username,
		Password:        config.DbConfig.Password,
		Database:        config.DbConfig.Database,
		ApplicationName: "demo",
		ReadTimeout:     ReadTimeout,
		WriteTimeout:    WriteTimeout,
		PoolSize:        PoolSize,
		MinIdleConns:    MinIdleConns,
	}
}

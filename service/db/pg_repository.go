package db

import "github.com/go-pg/pg/v9"

type PostgresClient interface {
	GetConnection() *pg.DB

	// Close closes the database client, releasing any open resources
	Close() error
}

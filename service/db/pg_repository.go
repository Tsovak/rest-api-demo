package db

import pg "github.com/go-pg/pg/v9"

//go:generate mockgen -source pg_repository.go -package mock -destination ../mock/pg_repository.go

// PostgresClient declare methods for Pg
type PostgresClient interface {
	GetConnection() *pg.DB

	// Close closes the database client, releasing any open resources
	Close() error
}

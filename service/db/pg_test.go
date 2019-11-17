package db

// +build integration

import (
	"github.com/go-pg/pg/v9"
	"github.com/stretchr/testify/require"
	"github.com/tsovak/rest-api-demo/config"
	"testing"
)

func TestPgClientWorking(t *testing.T) {
	cfg := config.Config{
		DbConfig: config.DbConfig{
			Address:    "localhost:5432",
			Username:   "test",
			Password:   "password",
			Sslmode:    "disable",
			Drivername: "postgres",
		},
	}
	client := NewPostgresClient(nil, cfg)
	require.NotNil(t, client)
	db := client.GetConnection()

	var num int
	_, err := db.Query(pg.Scan(&num), "SELECT ?", 42)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, 42, num)
}

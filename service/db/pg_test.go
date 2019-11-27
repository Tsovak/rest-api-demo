// +build integration

package db

import (
	"github.com/go-pg/pg/v9"
	"github.com/stretchr/testify/require"
	"github.com/tsovak/rest-api-demo/config"
	"github.com/tsovak/rest-api-demo/testutils"
	"testing"
)

func TestPgClientWorking(t *testing.T) {
	testConfig, err := config.GetTestConfig()
	require.NoError(t, err)

	setup, err := testutils.SetupTestDB(GetPgConnectionOptions(testConfig), "../../scripts/migrations/")
	require.NoError(t, err)

	client := NewPostgresClient(nil, setup.Db)
	require.NotNil(t, client)
	db := client.GetConnection()

	var num int
	_, err = db.Query(pg.Scan(&num), "SELECT ?", 42)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, 42, num)
}

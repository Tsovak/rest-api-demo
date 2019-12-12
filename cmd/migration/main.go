package main

import (
	"flag"
	"github.com/go-pg/migrations/v7"
	"github.com/pkg/errors"
	"github.com/tsovak/rest-api-demo/config"
	"github.com/tsovak/rest-api-demo/service/db"
	"os"
)

var migrationDir = flag.String("dir", "./scripts/migrations/", "directory with migrations")
var doInit = flag.Bool("init", true, "perform db init (for empty db)")

func main() {
	flag.Parse()

	config, err := config.LoadConfig()
	if err != nil {
		println(err.Error())
		os.Exit(-1)
	}

	pgClient := db.NewPostgresClientFromConfig(config)
	connection := pgClient.GetConnection()
	defer connection.Close()

	migrationCollection := migrations.NewCollection()
	if *doInit {
		_, _, err := migrationCollection.Run(connection, "init")
		if err != nil {
			config.Logger.Fatal(errors.Wrap(err, "Could not init migrations"))
		}
	}

	err = migrationCollection.DiscoverSQLMigrations(*migrationDir)
	if err != nil {
		config.Logger.Fatal(errors.Wrap(err, "Failed to read migrations"))
	}

	_, _, err = migrationCollection.Run(connection, "up")
	if err != nil {
		config.Logger.Fatal(errors.Wrap(err, "Could not migrate"))
	}
	config.Logger.Info("migrated successfully!")
}

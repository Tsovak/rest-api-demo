package config

var cfg = Config{
	DbConfig: DbConfig{
		Address:  "localhost",
		Username: "postgres",
		Password: "secret",
		//Database:   "test_db",
		Sslmode:    "disable",
		Drivername: "postgres",
	},
}

// GetTestConfig returns config for tests
func GetTestConfig() (Config, error) {
	return cfg, nil
}

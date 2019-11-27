package config

var cfg = Config{
	DbConfig: DbConfig{
		Address:  "localhost:5432",
		Username: "test",
		Password: "password",
		//Database:   "test",
		Sslmode:    "disable",
		Drivername: "postgres",
	},
}

func GetTestConfig() (Config, error) {
	return cfg, nil
}

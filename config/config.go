package config

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"strings"
)

const configFileName = "app"

type Config struct {
	ServerAddress string
	ServerPort    string
	DbConfig      DbConfig `mapstructure:"db"`
	Logger        *logrus.Logger
}

type DbConfig struct {
	Address    string
	Username   string
	Password   string
	Database   string
	Sslmode    string
	Drivername string
}

// LoadConfig load config from file
func LoadConfig() (Config, error) {
	v := viper.New()
	v.SetConfigName(configFileName)
	v.SetEnvPrefix("api")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.AddConfigPath("./config")

	replacer := strings.NewReplacer(".", "_")
	v.SetEnvKeyReplacer(replacer)
	v.AutomaticEnv()

	var cfg Config
	if err := v.ReadInConfig(); err != nil {
		return Config{}, errors.Wrap(err, "Failed to read config")
	}

	err := v.Unmarshal(&cfg)
	if err != nil {
		return Config{}, errors.Wrap(err, "Unable to decode into struct")
	}

	loglevel := v.GetString("loglevel")
	logger := getLogger(loglevel)

	cfg.Logger = logger
	return cfg, nil
}

func getLogger(loglevel string) *logrus.Logger {
	logrusLogger := logrus.New()
	logrusLogger.SetFormatter(&logrus.JSONFormatter{})
	level, ok := logrus.ParseLevel(loglevel)
	if ok != nil {
		level = logrus.DebugLevel
	}
	logrusLogger.SetLevel(level)
	return logrusLogger
}

// Load config without error
func LoadConfigSilence() Config {
	config, err := LoadConfig()
	if err != nil {
		panic(err)
	}
	return config
}

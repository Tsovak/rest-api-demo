package config

import (
	log "github.com/go-kit/kit/log/logrus"
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
	Logger        *logrus.FieldLogger
}

type DbConfig struct {
	Port       string
	Host       string
	Username   string
	Password   string
	Database   string
	Sslmode    string
	Drivername string
}

// LoadConfig load config from file
func LoadConfig() (Config, error) {
	viper.SetConfigName(configFileName)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.AutomaticEnv()

	var cfg Config
	if err := viper.ReadInConfig(); err != nil {
		return Config{}, errors.Wrap(err, "Failed to read config")
	}

	err := viper.Unmarshal(&cfg)
	if err != nil {
		return Config{}, errors.Wrap(err, "Unable to decode into struct")
	}

	loglevel := viper.GetString("loglevel")
	getLogger(loglevel)

	return cfg, nil
}

func getLogger(loglevel string) logrus.FieldLogger {
	logrusLogger := logrus.New()
	logrusLogger.SetFormatter(&logrus.JSONFormatter{})
	level, ok := logrus.ParseLevel(loglevel)
	if ok != nil {
		level = logrus.DebugLevel
	}
	logrusLogger.SetLevel(level)

	logger := log.NewLogrusLogger(logrusLogger).(logrus.FieldLogger)
	return logger
}

// Load config without error
func LoadConfigSilence() Config {
	config, err := LoadConfig()
	if err != nil {
		panic(err)
	}
	return config
}

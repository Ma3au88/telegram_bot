package config

import (
	"fmt"
	"github.com/spf13/viper"
	"time"
)

type Config struct {
	TelegramToken string
	Server
	Database
}

type Server struct {
	Port           string        `mapstructure:"port"`
	CurrencyUpdate time.Duration `mapstructure:"currencyUpdate"`
}

type Database struct {
	DBHost         string `mapstructure:"host"`
	DBPort         uint16 `mapstructure:"port"`
	DBUsername     string `mapstructure:"username"`
	DBName         string `mapstructure:"dbname"`
	DBPassword     string
	DBSSLmode      string `mapstructure:"sslmode"`
	DBPoolMaxConns uint16 `mapstructure:"pool_max_conns"`
}

func Init() (*Config, error) {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed reading configuration fales: %w", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal configuration: %w", err)
	}

	if err := viper.UnmarshalKey("server", &cfg.Server); err != nil {
		return nil, fmt.Errorf("unmarshal server configuration: %w", err)
	}

	if err := viper.UnmarshalKey("db", &cfg.Database); err != nil {
		return nil, fmt.Errorf("unmarshal database configuration: %w", err)
	}

	if err := parseEnv(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal env configuration: %w", err)
	}

	return &cfg, nil
}

func parseEnv(cfg *Config) error {
	if err := viper.BindEnv("TOKEN"); err != nil {
		return fmt.Errorf("error getting token from env: %w", err)
	}
	cfg.TelegramToken = viper.GetString("TOKEN")

	if err := viper.BindEnv("DB_PASSWORD"); err != nil {
		return fmt.Errorf("error getting database password from env: %w", err)
	}
	cfg.DBPassword = viper.GetString("DB_PASSWORD")

	return nil
}

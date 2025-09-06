package config

import (
	"github.com/spf13/viper"
	"time"
)

type Config struct {
	DBHost       string        `mapstructure:"postgres_host"`
	DBUser       string        `mapstructure:"postgres_user"`
	DBPass       string        `mapstructure:"postgres_password"`
	DBPort       int           `mapstructure:"postgres_port"`
	DBName       string        `mapstructure:"postgres_db"`
	SSLMode      string        `mapstructure:"postgres_ssl_mode"`
	JWTSecretKey string        `mapstructure:"jwt_secret_key"`
	TokenTTL     time.Duration `mapstructure:"token_ttl"`
	StartPort    string        `mapstructure:"start_port"`
}

func LoadConfig() (config Config, err error) {
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return config, nil
		}
	}

	err = viper.Unmarshal(&config)
	return
}

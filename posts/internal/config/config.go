package config

import (
	"errors"
	"github.com/spf13/viper"
	"time"
)

type Config struct {
	Postgres Postgres `mapstructure:"postgres"`
	HTTP     HTTP     `mapstructure:"http"`
	JWT      JWT      `mapstructure:"jwt"`
	Kafka    Kafka    `mapstructure:"kafka"`
}

type Postgres struct {
	Host     string `mapstructure:"host"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Port     string `mapstructure:"port"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
}

type HTTP struct {
	Port string `mapstructure:"port"`
}

type JWT struct {
	SecretKey string        `mapstructure:"secret_key"`
	TokenTTL  time.Duration `mapstructure:"token_ttl"`
}

type Kafka struct {
	Brokers []string `mapstructure:"brokers"`
	Topic   string   `mapstructure:"topic"`
}

func LoadConfig() (config Config, err error) {
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()
	if err = viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return config, errors.New("config file not found")
		}
		return config, err
	}
	err = viper.Unmarshal(&config)
	return
}

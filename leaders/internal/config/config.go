package config

import (
	"errors"
	"github.com/spf13/viper"
)

type Config struct {
	Env        string     `mapstructure:"env"`
	HTTPServer HTTPServer `mapstructure:"http_server"`
	Redis      Redis      `mapstructure:"redis"`
	Services   Services   `mapstructure:"services"`
	Kafka      Kafka      `mapstructure:"kafka"`
}

type Kafka struct {
	Brokers []string `mapstructure:"brokers"`
	Topic   string   `mapstructure:"topic"`
}

type HTTPServer struct {
	Port string `mapstructure:"port"`
}

type Redis struct {
	Addr string `mapstructure:"addr"`
}

type Services struct {
	UsersURL string `mapstructure:"users_url"`
}

func LoadConfig() (config Config, err error) {
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return config, errors.New("config file not found")
		}
	}

	err = viper.Unmarshal(&config)
	return
}

package main

import (
	"awesomeProject1/pkg/auth"
	"awesomeProject1/users/internal/config"
	"awesomeProject1/users/internal/handler"
	"awesomeProject1/users/internal/repository"
	"awesomeProject1/users/internal/service"
	"log"

	"github.com/sirupsen/logrus"
)

func main() {

	cfg, err := config.LoadConfig()
	if err != nil {
		logrus.Fatal(err)
	}

	db, err := repository.NewPostgresDB(repository.Config{
		Host:         cfg.DBHost,
		Port:         cfg.DBPort,
		Name:         cfg.DBUser,
		PasswordHash: cfg.DBPass,
		Database:     cfg.DBName,
	})
	if err != nil {
		logrus.Fatal(err)
	}

	tokenManager := auth.NewTokenManager([]byte(cfg.JWTSecretKey))

	repos := repository.NewRepository(db)
	services := service.NewAuthService(repos, []byte(cfg.JWTSecretKey), cfg.TokenTTL)
	handlers := handler.NewHandler(services, tokenManager)

	router := handlers.InitRoutes()

	log.Println("Starting server...")
	if err := router.Run(cfg.StartPort); err != nil {
		logrus.Fatal(err)
	}
}

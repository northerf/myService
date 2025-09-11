package main

import (
	"awesomeProject1/pkg/auth"
	"awesomeProject1/purchases/internal/config"
	"awesomeProject1/purchases/internal/handler"
	"awesomeProject1/purchases/internal/repository"
	"awesomeProject1/purchases/internal/service"
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
	services := service.NewPurchasesServices(repos)
	handlers := handler.NewHandler(services, tokenManager)

	router := handlers.InitRoutes()

	log.Println("Starting server...")
	if err := router.Run(cfg.StartPort); err != nil {
		logrus.Fatal(err)
	}
}

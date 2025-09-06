package main

import (
	"awesomeProject1/internal/config"
	"awesomeProject1/internal/handler"
	"awesomeProject1/internal/repository"
	"awesomeProject1/internal/service"
	"github.com/sirupsen/logrus"
	"log"
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

	repos := repository.NewRepository(db)
	services := service.NewAuthService(repos, []byte(cfg.JWTSecretKey), cfg.TokenTTL)
	handlers := handler.NewHandler(services)

	router := handlers.InitRoutes()

	log.Println("Starting server...")
	if err := router.Run(cfg.StartPort); err != nil {
		logrus.Fatal(err)
	}
}

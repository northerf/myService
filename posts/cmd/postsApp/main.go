package main

import (
	"awesomeProject1/posts/internal/config"
	"awesomeProject1/posts/internal/handler"
	"awesomeProject1/posts/internal/repository"
	"awesomeProject1/posts/internal/service"
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

	repos := repository.NewRepository(db)
	services := service.NewPostsService(repos)
	handlers := handler.NewHandler(services)

	router := handlers.InitRoutes()

	log.Println("Starting server...")
	if err := router.Run(cfg.StartPort); err != nil {
		logrus.Fatal(err)
	}
}

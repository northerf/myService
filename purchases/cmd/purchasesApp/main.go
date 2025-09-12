package main

import (
	"awesomeProject1/pkg/auth"
	"awesomeProject1/purchases/internal/config"
	"awesomeProject1/purchases/internal/handler"
	"awesomeProject1/purchases/internal/kafka"
	"awesomeProject1/purchases/internal/repository"
	"awesomeProject1/purchases/internal/service"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	port, err := strconv.Atoi(cfg.Postgres.Port)
	if err != nil {
		log.Fatalf("invalid postgres port: %v", err)
	}

	db, err := repository.NewPostgresDB(repository.Config{
		Host:         cfg.Postgres.Host,
		Port:         port,
		Name:         cfg.Postgres.User,
		PasswordHash: cfg.Postgres.Password,
		Database:     cfg.Postgres.DBName,
	})
	if err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}

	tokenManager := auth.NewTokenManager([]byte(cfg.JWT.SecretKey))

	repos := repository.NewRepository(db)
	kafkaProducer := kafka.NewProducer(cfg.Kafka.Brokers, cfg.Kafka.Topic)
	defer kafkaProducer.Close()

	services := service.NewService(repos, kafkaProducer)
	handlers := handler.NewHandler(services, tokenManager)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.HTTP.Port),
		Handler: handlers.InitRoutes(),
	}

	go func() {
		log.Printf("Starting server on port %s", cfg.HTTP.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start HTTP server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}

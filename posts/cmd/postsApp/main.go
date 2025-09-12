package main

import (
	"awesomeProject1/posts/internal/config"
	"awesomeProject1/posts/internal/handler"
	"awesomeProject1/posts/internal/kafka"
	"awesomeProject1/posts/internal/repository"
	"awesomeProject1/posts/internal/service"
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

	repos := repository.NewRepository(db)
	services := service.NewPostsService(repos)
	handlers := handler.NewHandler(services)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if len(cfg.Kafka.Brokers) > 0 && cfg.Kafka.Topic != "" {
		kafkaConsumer := kafka.NewConsumer(cfg.Kafka.Brokers, cfg.Kafka.Topic, services)
		go func() {
			log.Printf("Starting Kafka consumer: brokers=%v topic=%s", cfg.Kafka.Brokers, cfg.Kafka.Topic)
			kafkaConsumer.Run(ctx)
		}()
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.HTTP.Port),
		Handler: handlers.InitRoutes(),
	}

	go func() {
		log.Printf("Starting HTTP server on port %s", cfg.HTTP.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start HTTP server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	cancel()

	ctxShutdown, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutdown()

	if err := srv.Shutdown(ctxShutdown); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}

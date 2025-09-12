package main

import (
	"awesomeProject1/leaders/internal/config"
	"awesomeProject1/leaders/internal/handler"
	"awesomeProject1/leaders/internal/kafka"
	redis_repository "awesomeProject1/leaders/internal/redis-repository"
	"awesomeProject1/leaders/internal/service"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/redis/go-redis/v9"
)

func main() {
	log.Println("Loading configuration...")
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Config loaded: env=%s http_port=%s redis_addr=%s users_url=%s kafka_brokers=%v kafka_topic=%s", cfg.Env, cfg.HTTPServer.Port, cfg.Redis.Addr, cfg.Services.UsersURL, cfg.Kafka.Brokers, cfg.Kafka.Topic)

	log.Println("Connecting to Redis...")
	redisClient := redis.NewClient(&redis.Options{
		Addr: cfg.Redis.Addr,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Fatalf("Cannot connect to Redis: %v", err)
	}
	log.Println("Successfully connected to Redis")

	repo := redis_repository.NewRedisRepository(redisClient)
	userProvider := service.NewUserServiceClient(cfg.Services.UsersURL)
	leadersService := service.NewLeadersService(repo, userProvider)

	h := handler.NewHandler(leadersService)

	rootCtx, rootCancel := context.WithCancel(context.Background())
	defer rootCancel()

	g, gCtx := errgroup.WithContext(rootCtx)

	g.Go(func() error {
		if len(cfg.Kafka.Brokers) == 0 || cfg.Kafka.Topic == "" {
			log.Println("Kafka config missing (brokers/topic). Consumer is disabled.")
			<-gCtx.Done()
			return nil
		}
		log.Printf("Starting Kafka consumer: brokers=%v topic=%s", cfg.Kafka.Brokers, cfg.Kafka.Topic)
		consumer := kafka.NewConsumer(cfg.Kafka.Brokers, cfg.Kafka.Topic, leadersService)
		consumer.Run(gCtx)
		return nil
	})

	router := h.InitRoutes()
	for _, r := range router.Routes() {
		log.Printf("Route registered: %s %s -> %s", r.Method, r.Path, r.Handler)
	}

	srv := &http.Server{
		Addr:    cfg.HTTPServer.Port,
		Handler: router,
	}

	go func() {
		log.Printf("Starting HTTP server on port %s", cfg.HTTPServer.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start HTTP server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-gCtx.Done():
		log.Println("Context cancelled")
	case sig := <-quit:
		log.Println("Got signal:", sig)
	}

	log.Println("Initiating graceful shutdown...")
	rootCancel()

	ctxShutdown, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutdown()

	if err := srv.Shutdown(ctxShutdown); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	done := make(chan struct{})
	go func() {
		defer close(done)
		if err := g.Wait(); err != nil && err != http.ErrServerClosed {
			log.Printf("errgroup error: %v\n", err)
		}
	}()
	select {
	case <-done:
		log.Println("Background workers stopped")
	case <-time.After(5 * time.Second):
		log.Println("Timeout waiting for background workers; exiting")
	}

	log.Println("Server exiting")
}

package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/ashucfx/distributed-job-platform/packages/config"
	"github.com/ashucfx/distributed-job-platform/packages/database"
	"github.com/ashucfx/distributed-job-platform/packages/logger"
	"github.com/ashucfx/distributed-job-platform/packages/queue"
)

func main() {
	cfg := config.LoadConfig()
	logger.InitLogger(cfg.AppEnv)
	defer logger.Sync()

	dbConn := database.Connect(cfg.DBUrl)
	dbService := database.NewGormDBService(dbConn)
	queueService := queue.NewRedisQueue(cfg.RedisUrl)

	engine := NewEngine(dbService, queueService)
	
	// Register handlers
	engine.Register("send_email", &SendEmailHandler{})

	// Context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle sigterm
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	
	go func() {
		<-signals
		logger.Log.Info("received shutdown signal, stopping workers...")
		cancel()
	}()

	// Run with 5 concurrent workers
	engine.Run(ctx, 5)
}

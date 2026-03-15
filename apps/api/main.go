package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/ashucfx/distributed-job-platform/packages/config"
	"github.com/ashucfx/distributed-job-platform/packages/database"
	"github.com/ashucfx/distributed-job-platform/packages/logger"
	"github.com/ashucfx/distributed-job-platform/packages/queue"
)

func main() {
	cfg := config.LoadConfig()
	logger.InitLogger(cfg.AppEnv)
	defer logger.Sync()

	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	dbConn := database.Connect(cfg.DBUrl)
	dbService := database.NewGormDBService(dbConn)
	queueService := queue.NewRedisQueue(cfg.RedisUrl)

	jobController := NewJobController(dbService, queueService)
	router := SetupRouter(jobController)

	log.Printf("Starting API server on port %s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
}

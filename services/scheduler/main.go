package main

import (
	"context"
	"encoding/json"
	"os"
	"os/signal"
	"syscall"

	"github.com/google/uuid"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"

	"github.com/ashucfx/distributed-job-platform/packages/config"
	"github.com/ashucfx/distributed-job-platform/packages/database"
	"github.com/ashucfx/distributed-job-platform/packages/logger"
	"github.com/ashucfx/distributed-job-platform/packages/queue"
)

type SchedulerEngine struct {
	db   database.DBService
	q    queue.QueueService
	cron *cron.Cron
}

func NewSchedulerEngine(db database.DBService, q queue.QueueService) *SchedulerEngine {
	return &SchedulerEngine{
		db:   db,
		q:    q,
		cron: cron.New(cron.WithSeconds()), // High precision for demo
	}
}

func (s *SchedulerEngine) AddSystemJob(schedule string, name string, payload map[string]interface{}) {
	_, err := s.cron.AddFunc(schedule, func() {
		logger.Log.Info("triggering scheduled job", zap.String("name", name))

		payloadBytes, _ := json.Marshal(payload)
		
		job := database.Job{
			ID:      uuid.New().String(),
			Name:    name,
			Payload: string(payloadBytes),
			Status:  database.StatusPending,
		}

		if err := s.db.CreateJob(&job); err != nil {
			logger.Log.Error("failed to create scheduled job in db", zap.Error(err))
			return
		}

		err := s.q.Enqueue(context.Background(), queue.JobPayload{
			ID:      job.ID,
			Name:    job.Name,
			Payload: job.Payload,
		})

		if err != nil {
			logger.Log.Error("failed to enqueue scheduled job", zap.Error(err))
		}
	})

	if err != nil {
		logger.Log.Error("failed to register cron job", zap.Error(err))
	}
}

func (s *SchedulerEngine) Start() {
	s.cron.Start()
	logger.Log.Info("scheduler started")
}

func (s *SchedulerEngine) Stop() {
	<-s.cron.Stop().Done()
	logger.Log.Info("scheduler stopped")
}

func main() {
	cfg := config.LoadConfig()
	logger.InitLogger(cfg.AppEnv)
	defer logger.Sync()

	dbConn := database.Connect(cfg.DBUrl)
	dbService := database.NewGormDBService(dbConn)
	queueService := queue.NewRedisQueue(cfg.RedisUrl)

	engine := NewSchedulerEngine(dbService, queueService)

	// Register some default cron jobs
	engine.AddSystemJob("*/30 * * * * *", "send_email", map[string]interface{}{
		"to": "admin@example.com",
		"subject": "System Health Check (Every 30s)",
	})

	engine.Start()

	// Wait for terminate signal
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	<-signals

	logger.Log.Info("shutting down scheduler engine")
	engine.Stop()
}

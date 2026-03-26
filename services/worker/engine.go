package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/ashucfx/distributed-job-platform/packages/database"
	"github.com/ashucfx/distributed-job-platform/packages/logger"
	"github.com/ashucfx/distributed-job-platform/packages/queue"
)

type JobHandler interface {
	Handle(ctx context.Context, payload map[string]interface{}) error
}

type Engine struct {
	db       database.DBService
	q        queue.QueueService
	handlers map[string]JobHandler
}

func NewEngine(db database.DBService, q queue.QueueService) *Engine {
	return &Engine{
		db:       db,
		q:        q,
		handlers: make(map[string]JobHandler),
	}
}

func (e *Engine) Register(name string, handler JobHandler) {
	e.handlers[name] = handler
}

func (e *Engine) Run(ctx context.Context, concurrency int) {
	logger.Log.Info("starting worker engine", zap.Int("concurrency", concurrency))

	for i := 0; i < concurrency; i++ {
		go e.workerLoop(ctx, i)
	}

	<-ctx.Done()
	logger.Log.Info("shutting down worker engine")
}

func (e *Engine) workerLoop(ctx context.Context, id int) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			job, err := e.q.Dequeue(ctx, 2*time.Second)
			if err != nil {
				logger.Log.Error("failed to dequeue job", zap.Error(err), zap.Int("worker_id", id))
				time.Sleep(1 * time.Second)
				continue
			}

			if job == nil {
				continue // Timeout, no jobs
			}

			e.processJob(ctx, job, id)
		}
	}
}

func (e *Engine) processJob(ctx context.Context, jobPayload *queue.JobPayload, workerID int) {
	logger.Log.Info("processing job", zap.String("job_id", jobPayload.ID), zap.Int("worker_id", workerID))

	// Mark as processing
	err := e.db.UpdateJobStatus(jobPayload.ID, database.StatusProcessing, "")
	if err != nil {
		logger.Log.Error("failed to update job status to processing", zap.Error(err))
		return
	}
	_ = e.db.CreateJobLog(jobPayload.ID, fmt.Sprintf("Worker %d started processing", workerID), "info")

	handler, exists := e.handlers[jobPayload.Name]
	if !exists {
		errMsg := fmt.Sprintf("no handler registered for job type: %s", jobPayload.Name)
		e.handleFailure(ctx, jobPayload.ID, jobPayload, errMsg)
		return
	}

	var payloadMap map[string]interface{}
	if err := json.Unmarshal([]byte(jobPayload.Payload), &payloadMap); err != nil {
		e.handleFailure(ctx, jobPayload.ID, jobPayload, "invalid payload format")
		return
	}

	// Execution
	execErr := handler.Handle(ctx, payloadMap)
	
	if execErr != nil {
		e.handleFailure(ctx, jobPayload.ID, jobPayload, execErr.Error())
	} else {
		logger.Log.Info("job completed successfully", zap.String("job_id", jobPayload.ID))
		_ = e.db.UpdateJobStatus(jobPayload.ID, database.StatusCompleted, "")
		_ = e.db.CreateJobLog(jobPayload.ID, "Job completed successfully", "info")
	}
}

func (e *Engine) handleFailure(ctx context.Context, dbJobID string, job *queue.JobPayload, errMsg string) {
	logger.Log.Error("job failed", zap.String("job_id", job.ID), zap.String("error", errMsg))
	
	dbJob, err := e.db.GetJobByID(dbJobID)
	if err != nil {
		logger.Log.Error("failed to retrieve job logic on failure", zap.Error(err))
		return
	}

	_ = e.db.CreateJobLog(dbJobID, "Job failed: "+errMsg, "error")

	// Retry logic
	if dbJob.Retries < dbJob.MaxRetries {
		_ = e.db.IncrementRetries(dbJobID)
		
		// In a production system, we'd use a dedicated delayed queue or a scheduler.
		// For this platform, we re-enqueue with a backoff.
		backoffDuration := time.Duration(1<<dbJob.Retries) * time.Minute
		_ = e.db.CreateJobLog(dbJobID, fmt.Sprintf("Retrying job in %v (Retry %d/%d)", backoffDuration, dbJob.Retries+1, dbJob.MaxRetries), "info")
		
		_ = e.db.UpdateJobStatus(dbJobID, database.StatusPending, errMsg)

		go func() {
			time.Sleep(backoffDuration)
			_ = e.q.Enqueue(context.Background(), *job)
		}()
	} else {
		// DLQ
		_ = e.db.CreateJobLog(dbJobID, "Max retries reached. Moving to DLQ.", "error")
		err := e.q.EnqueueDLQ(ctx, *job)
		if err != nil {
			logger.Log.Error("failed to move to DLQ", zap.Error(err))
		}
		_ = e.db.UpdateJobStatus(dbJobID, database.StatusFailed, errMsg)
	}
}

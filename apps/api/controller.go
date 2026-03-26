package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/ashucfx/distributed-job-platform/packages/database"
	"github.com/ashucfx/distributed-job-platform/packages/logger"
	"github.com/ashucfx/distributed-job-platform/packages/queue"
)

type JobController struct {
	db    database.DBService
	queue queue.QueueService
}

func NewJobController(db database.DBService, q queue.QueueService) *JobController {
	return &JobController{db: db, queue: q}
}

type CreateJobRequest struct {
	Name    string                 `json:"name" binding:"required"`
	Payload map[string]interface{} `json:"payload"`
}

func (c *JobController) CreateJob(ctx *gin.Context) {
	var req CreateJobRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	payloadBytes, _ := json.Marshal(req.Payload)
	job := database.Job{
		ID:      uuid.New().String(),
		Name:    req.Name,
		Payload: string(payloadBytes),
		Status:  database.StatusPending,
	}

	if err := c.db.CreateJob(&job); err != nil {
		logger.Log.Error("failed to create job in db", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create job"})
		return
	}

	err := c.queue.Enqueue(context.Background(), queue.JobPayload{
		ID:      job.ID,
		Name:    job.Name,
		Payload: job.Payload,
	})

	if err != nil {
		logger.Log.Error("failed to enqueue job", zap.Error(err))
		// Note: At a real scale, we'd want a transaction outbox pattern here.
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to enqueue job"})
		return
	}

	ctx.JSON(http.StatusCreated, job)
}

func (c *JobController) GetJob(ctx *gin.Context) {
	id := ctx.Param("id")
	job, err := c.db.GetJobByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "job not found"})
		return
	}
	ctx.JSON(http.StatusOK, job)
}

func (c *JobController) ListJobs(ctx *gin.Context) {
	jobs, err := c.db.ListJobs(20)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list jobs"})
		return
	}
	ctx.JSON(http.StatusOK, jobs)
}

package queue

import (
	"context"
	"time"
)

type QueueService interface {
	Enqueue(ctx context.Context, job JobPayload) error
	Dequeue(ctx context.Context, timeout time.Duration) (*JobPayload, error)
	EnqueueDLQ(ctx context.Context, job JobPayload) error
}

package queue

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	QueueName   = "jobs:queue"
	DeadLetterQ = "jobs:dlq"
)

type RedisQueue struct {
	client *redis.Client
}

func NewRedisQueue(redisURL string) *RedisQueue {
	client := redis.NewClient(&redis.Options{
		Addr: redisURL,
	})
	return &RedisQueue{client: client}
}

func (r *RedisQueue) Enqueue(ctx context.Context, job JobPayload) error {
	data, err := json.Marshal(job)
	if err != nil {
		return err
	}
	return r.client.LPush(ctx, QueueName, data).Err()
}

func (r *RedisQueue) Dequeue(ctx context.Context, timeout time.Duration) (*JobPayload, error) {
	result, err := r.client.BRPop(ctx, timeout, QueueName).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Timeout, no job
		}
		return nil, err
	}

	// result[0] is list name, result[1] is the value
	var job JobPayload
	if err := json.Unmarshal([]byte(result[1]), &job); err != nil {
		return nil, err
	}
	return &job, nil
}

func (r *RedisQueue) EnqueueDLQ(ctx context.Context, job JobPayload) error {
	data, err := json.Marshal(job)
	if err != nil {
		return err
	}
	return r.client.LPush(ctx, DeadLetterQ, data).Err()
}

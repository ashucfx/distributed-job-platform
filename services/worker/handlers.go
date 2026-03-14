package main

import (
	"context"
	"fmt"
	"time"
)

type SendEmailHandler struct{}

func (h *SendEmailHandler) Handle(ctx context.Context, payload map[string]interface{}) error {
	to, ok := payload["to"].(string)
	if !ok {
		return fmt.Errorf("missing 'to' in payload")
	}
	subject, ok := payload["subject"].(string)
	if !ok {
		return fmt.Errorf("missing 'subject' in payload")
	}

	// Simulate work
	fmt.Printf("Simulating sending email to %s with subject: %s\n", to, subject)
	time.Sleep(2 * time.Second)
	
	// Simulate random failure for retry demonstration
	if time.Now().UnixNano()%3 == 0 {
		return fmt.Errorf("simulated temporary network error during email send")
	}

	fmt.Printf("Successfully sent email to %s\n", to)
	return nil
}

package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/hibiken/asynq"
)

type Consumer struct {
	queue  string
	client *asynq.Server
	server *asynq.ServeMux
}

func NewConsumer(queue string) *Consumer {
	redisAddr := os.Getenv("REDIS_ADDRESS")
	clientConsumer := asynq.NewServer(
		asynq.RedisClientOpt{Addr: redisAddr},
		asynq.Config{
			Concurrency: 1,
		},
	)

	mux := asynq.NewServeMux()

	return &Consumer{
		queue:  queue,
		client: clientConsumer,
		server: mux,
	}

}

func (c *Consumer) Consumer(handler func(message map[string]interface{}) error) error {
	c.server.HandleFunc(c.queue, func(ctx context.Context, task *asynq.Task) error {
		var item map[string]interface{}
		if err := json.Unmarshal(task.Payload(), &item); err != nil {
			fmt.Errorf("json.Unmarshal failed: %v", err)
			return err
		}

		err := handler(item)
		if err != nil {
			return err
		}

		return nil
	})

	return nil
}

func (c *Consumer) Start() {
	if err := c.client.Run(c.server); err != nil {
		log.Fatalf("could not run server: %v", err)
	}

}

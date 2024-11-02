package queue

import (
	"encoding/json"
	"os"
	"time"

	"github.com/hibiken/asynq"
)

type Publisher struct {
	queue  string
	client *asynq.Client
}

func NewPublisher(queue string) *Publisher {
	redisAddr := os.Getenv("REDIS_ADDRESS")
	client := asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddr})
	return &Publisher{
		queue:  queue,
		client: client,
	}
}

func (p *Publisher) Publish(message interface{}, totalRetries int) error {
	payloadMessage, _ := json.Marshal(message)
	_, error := p.client.Enqueue(
		asynq.NewTask(p.queue, []byte(payloadMessage)),
	)

	if error != nil {
		return error
	}

	return nil
}

func (p *Publisher) PublishWithDelay(
	message interface{},
	totalRetries int,
	delay time.Duration,
) error {
	payloadMessage, _ := json.Marshal(message)
	_, error := p.client.Enqueue(
		asynq.NewTask(p.queue, []byte(payloadMessage)),
		asynq.ProcessIn(delay),
	)

	if error != nil {
		return error
	}

	return nil
}

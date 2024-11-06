package queue

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/nats-io/stan.go"
)

type Publisher struct {
	queue  string
	client stan.Conn
}

func NewPublisher(queue string) *Publisher {
	natsURL := os.Getenv("NATS_ADDRESS")
	opts := []stan.Option{
		stan.NatsURL(natsURL),
	}

	conn, err := stan.Connect("test-cluster", uuid.NewString(), opts...)
	if err != nil {
		log.Fatal("Failed to connect to NATS", err)
	}

	return &Publisher{
		queue:  queue,
		client: conn,
	}
}

func (p *Publisher) Publish(message interface{}, totalRetries int) error {
	payloadMessage, _ := json.Marshal(message)
	error := p.client.Publish(
		p.queue, payloadMessage,
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
	ticker := time.NewTicker(delay)
	go func(ticker *time.Ticker) {
		for {
			select {
			case <-ticker.C:
				payloadMessage, _ := json.Marshal(message)
				err := p.client.Publish(
					p.queue, payloadMessage,
				)

				if err != nil {
					fmt.Sprintf("Err: %v", err)
				}

				ticker.Stop()
			}
		}
	}(ticker)
	return nil
}

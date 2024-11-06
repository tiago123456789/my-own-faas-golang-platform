package queue

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/google/uuid"
	"github.com/nats-io/stan.go"
)

type Consumer struct {
	queue   string
	client  stan.Conn
	handler func(message map[string]interface{}) error
}

func NewConsumer(queue string, handler func(message map[string]interface{}) error) *Consumer {
	natsURL := os.Getenv("NATS_ADDRESS")
	opts := []stan.Option{
		stan.NatsURL(natsURL),
	}

	conn, err := stan.Connect("test-cluster", uuid.NewString(), opts...)
	if err != nil {
		log.Fatal("Failed to connect to NATS", err)
	}

	return &Consumer{
		queue:   queue,
		client:  conn,
		handler: handler,
	}

}

func (c *Consumer) Start() {
	_, err := c.client.QueueSubscribe(c.queue, c.queue, func(msg *stan.Msg) {
		var item map[string]interface{}
		if err := json.Unmarshal(msg.Data, &item); err != nil {
			fmt.Errorf("json.Unmarshal failed: %v", err)
			return
		}

		err := c.handler(item)
		if err != nil {
			fmt.Errorf("json.Unmarshal failed: %v", err)
			return
		}

		msg.Ack()
	},
		stan.DurableName("durable-sub"),
		stan.SetManualAckMode(),
		stan.DeliverAllAvailable(),
	)

	if err != nil {
		fmt.Errorf("json.Unmarshal failed: %v", err)
	}

	runtime.Goexit()
}

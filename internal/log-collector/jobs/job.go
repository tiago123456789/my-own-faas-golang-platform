package jobs

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/elastic/go-elasticsearch/v7"
	esapi "github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/tiago123456789/my-own-faas-golang-platform/internal/log-collector/types"
	"github.com/tiago123456789/my-own-faas-golang-platform/pkg/queue"
)

func Init(es *elasticsearch.Client) {

	consumer := queue.NewConsumer(
		"logs",
		func(message map[string]interface{}) error {
			fmt.Println("starting to process message")
			var logData types.Log

			jsonString, _ := json.Marshal(message)
			fmt.Println(string(jsonString))

			err := json.Unmarshal(jsonString, &logData)
			if err != nil {
				fmt.Println(
					fmt.Sprintf("Error: %v", err),
				)
				return err
			}

			logData.Level = strings.ToLower(logData.Level)
			logToByte, err := json.Marshal(logData)
			if err != nil {
				fmt.Println(
					fmt.Sprintf("Error: %v", err),
				)
				return err
			}

			req := esapi.IndexRequest{
				Index:   "logs",
				Body:    bytes.NewReader(logToByte),
				Refresh: "true",
			}

			res, err := req.Do(context.Background(), es)
			if err != nil {
				log.Fatalf("Error getting response: %s", err)
			}
			defer res.Body.Close()

			if res.IsError() {
				log.Printf("[%s] Error indexing document ID", res.Status())
			}
			return nil
		},
	)

	consumer.Start()
}

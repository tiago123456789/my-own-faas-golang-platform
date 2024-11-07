package repositories

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/tiago123456789/my-own-faas-golang-platform/internal/log-collector/types"
)

type LogRepository struct {
	esClient *elasticsearch.Client
}

func NewLogRepository(
	esClient *elasticsearch.Client,
) *LogRepository {
	return &LogRepository{
		esClient: esClient,
	}
}

func (r *LogRepository) Register(message map[string]interface{}) error {
	var logData types.Log

	jsonString, _ := json.Marshal(message)

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

	res, err := req.Do(context.Background(), r.esClient)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Printf("[%s] Error indexing document ID", res.Status())
		return errors.New(fmt.Sprint("[%s] Error indexing document ID", res.Status()))
	}

	return nil
}

package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/tiago123456789/my-own-faas-golang-platform/internal/faas-manager/models"
	"github.com/tiago123456789/my-own-faas-golang-platform/internal/faas-manager/types"
	"github.com/tiago123456789/my-own-faas-golang-platform/pkg/queue"
	"gorm.io/gorm"
)

type FunctionService struct {
	db        *gorm.DB
	esClient  *elasticsearch.Client
	publisher queue.Publisher
}

func NewFunctionService(
	db *gorm.DB,
	publisher queue.Publisher,
	esClient *elasticsearch.Client,
) *FunctionService {
	return &FunctionService{
		db:        db,
		publisher: publisher,
		esClient:  esClient,
	}
}

func (f *FunctionService) GetLogs(functionName string) []types.Log {
	query := fmt.Sprintf(`{
		"query": {
			"bool": {
				"must": [
					{
						"term": {
							"service.keyword": "%s"
						}
					}
				]
			}
		},
		"sort": [
			{
				"timestamp": {
					"order": "desc"
				}
			}
		],
		"size": 100
	}`, functionName)

	res, err := f.esClient.Search(
		f.esClient.Search.WithContext(context.Background()),
		f.esClient.Search.WithIndex("logs"),
		f.esClient.Search.WithBody(strings.NewReader(query)),
		f.esClient.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		log.Fatalf("Error executing search request: %s", err)
	}
	defer res.Body.Close()

	var response struct {
		Hits struct {
			Hits []struct {
				Source types.Log `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}

	var logs []types.Log
	for _, hit := range response.Hits.Hits {
		logs = append(logs, hit.Source)
	}
	return logs
}

func (f *FunctionService) FindById(id string) models.Function {
	var function models.Function
	f.db.First(&function, "id = ?", id)
	return function
}

func (f *FunctionService) FindAll() []models.Function {
	var functions []models.Function
	f.db.Find(&functions)
	return functions
}

func (f *FunctionService) Deploy(newFunction types.NewFunction, lambdaPath string) (int, error) {
	function := models.Function{
		LambdaName:    newFunction.Name,
		Runtime:       newFunction.Runtime,
		LambdaPath:    lambdaPath,
		Cpu:           newFunction.Cpu,
		Memory:        newFunction.Memory,
		BuildProgress: "PENDENT",
	}

	var functionReturned models.Function
	f.db.First(&functionReturned, "lambda_name = ?", newFunction.Name)
	if functionReturned.ID == 0 {
		f.db.Create(&function)
	} else {
		function.ID = functionReturned.ID
		f.db.Updates(function)
	}

	newFunction.ID = function.ID
	newFunction.LambdaPath = lambdaPath
	err := f.publisher.Publish(newFunction, 2)
	if err != nil {
		return 0, err
	}

	return function.ID, nil
}

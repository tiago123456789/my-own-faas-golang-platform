package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/tiago123456789/my-own-faas-golang-platform/internal/faas-manager/models"
	"github.com/tiago123456789/my-own-faas-golang-platform/internal/faas-manager/types"
	"gorm.io/gorm"
)

type FunctionRepository struct {
	db       *gorm.DB
	esClient *elasticsearch.Client
}

func NewFunctionRepository(
	db *gorm.DB,
	esClient *elasticsearch.Client,
) *FunctionRepository {
	return &FunctionRepository{
		db:       db,
		esClient: esClient,
	}
}

func (r *FunctionRepository) FindById(id string) models.Function {
	var function models.Function
	r.db.First(&function, "id = ?", id)
	return function
}

func (r *FunctionRepository) Create(newFunction *models.Function) {
	r.db.Create(&newFunction)
}

func (r *FunctionRepository) Update(functionModified models.Function) {
	r.db.Updates(functionModified)
}

func (r *FunctionRepository) FindByName(name string) models.Function {
	var function models.Function
	r.db.First(&function, "lambda_name = ?", name)
	return function
}

func (r *FunctionRepository) FindAll() []models.Function {
	var functions []models.Function
	r.db.Find(&functions)
	return functions
}

func (r *FunctionRepository) GetLogs(functionName string) []types.Log {
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

	res, err := r.esClient.Search(
		r.esClient.Search.WithContext(context.Background()),
		r.esClient.Search.WithIndex("logs"),
		r.esClient.Search.WithBody(strings.NewReader(query)),
		r.esClient.Search.WithTrackTotalHits(true),
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

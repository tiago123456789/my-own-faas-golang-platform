package configs

import (
	"os"

	"github.com/elastic/go-elasticsearch/v7"
)

func InitES() (*elasticsearch.Client, error) {
	elasticsearchURL := os.Getenv("ELASTICSEARCH_URL")
	cfg := elasticsearch.Config{
		Addresses: []string{elasticsearchURL},
	}

	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	return es, nil
}

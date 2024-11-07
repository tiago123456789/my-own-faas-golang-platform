package jobs

import (
	"github.com/tiago123456789/my-own-faas-golang-platform/internal/log-collector/repositories"
	"github.com/tiago123456789/my-own-faas-golang-platform/pkg/queue"
)

func Init(repository *repositories.LogRepository) {

	consumer := queue.NewConsumer(
		"logs",
		func(message map[string]interface{}) error {
			err := repository.Register(message)
			return err
		},
	)

	consumer.Start()
}

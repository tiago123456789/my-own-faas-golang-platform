package lambda

import (
	"fmt"

	"go.uber.org/zap"
)

func Handler(logger *zap.Logger) error {
	logger.Info("Starting code cronjob hello-world-cron2")

	fmt.Println("HELLO WORLD 2!!!!")

	logger.Info("Finished code cronjob hello-world-cron2")

	return nil
}

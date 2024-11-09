package lambda

import (
	"fmt"

	"go.uber.org/zap"
)

func Handler(logger *zap.Logger) error {
	logger.Info("Starting code cronjob")

	fmt.Println("HELLO WORLD!!!!")

	return nil
}

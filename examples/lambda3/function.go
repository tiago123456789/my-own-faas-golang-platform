package lambda3

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func Handler(c *fiber.Ctx, logger *zap.Logger) error {
	// logger.Info("Execute lambda function in golang")
	return c.JSON(fiber.Map{
		"id":      uuid.NewString(),
		"message": os.Getenv("MESSAGE"),
	})
}

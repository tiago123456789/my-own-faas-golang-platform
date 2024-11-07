package lambda

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func Handler(c *fiber.Ctx, logger *zap.Logger) error {
	token := c.Params("token")
	logger.Info("Execute lambda function in golang")
	return c.JSON(fiber.Map{
		"token":       token,
		"id":          uuid.NewString(),
		"message":     "Hi world!!!!",
		"env_message": os.Getenv("MESSAGE"),
	})
}

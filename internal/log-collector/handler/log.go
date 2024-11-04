package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/tiago123456789/my-own-faas-golang-platform/internal/log-collector/types"
	"github.com/tiago123456789/my-own-faas-golang-platform/pkg/queue"
)

type LogHandler struct {
	logPublisher queue.Publisher
}

func NewLogHandler(
	logPublisher queue.Publisher,
) *LogHandler {
	return &LogHandler{
		logPublisher: logPublisher,
	}
}

func (l *LogHandler) Register(c *fiber.Ctx) error {
	logData := new(types.Log)
	if err := c.BodyParser(logData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	err := l.logPublisher.Publish(logData, 2)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.SendStatus(202)
}

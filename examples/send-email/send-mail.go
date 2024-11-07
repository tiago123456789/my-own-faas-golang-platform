package lambda

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	gomail "gopkg.in/mail.v2"
)

func Handler(c *fiber.Ctx, logger *zap.Logger) error {
	message := gomail.NewMessage()

	message.SetHeader("From", "youremail@email.com")
	message.SetHeader("To", "recipient1@email.com")

	subject := os.Getenv("SUBJECT")
	messageEmail := os.Getenv("MESSAGE")
	message.SetHeader("Subject", subject)

	message.SetBody("text/plain", messageEmail)

	dialer := gomail.NewDialer("sandbox.smtp.mailtrap.io", 587, "user_here", "password_here")

	if err := dialer.DialAndSend(message); err != nil {
		fmt.Println("Error:", err)
		return c.Status(500).JSON(fiber.Map{
			"message": err.Error(),
		})
	} else {
		return c.Status(200).JSON(fiber.Map{
			"message": "Email sent successfully!",
		})
	}
}

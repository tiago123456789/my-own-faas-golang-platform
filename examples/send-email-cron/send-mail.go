package lambda

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	gomail "gopkg.in/mail.v2"
)

func Handler(logger *zap.Logger) error {
	message := gomail.NewMessage()

	message.SetHeader("From", "youremail@email.com")
	message.SetHeader("To", "recipient1@email.com")

	subject := os.Getenv("SUBJECT")
	messageEmail := os.Getenv("MESSAGE")
	message.SetHeader("Subject", subject)

	message.SetBody("text/plain", messageEmail)

	dialer := gomail.NewDialer("sandbox.smtp.mailtrap.io", 587, "", "")

	if err := dialer.DialAndSend(message); err != nil {
		fmt.Println("Error:", err)
	}

	return nil
}

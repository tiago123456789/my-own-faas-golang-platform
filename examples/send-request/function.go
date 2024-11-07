package lambda

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type Data struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func Handler(c *fiber.Ctx, logger *zap.Logger) error {
	logger.Info("Starting process to send webhook notification")
	url := os.Getenv("WEBHOOK_URL")

	var payload Data

	if err := c.BodyParser(&payload); err != nil {
		return err
	}

	// Marshal the data into JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return c.SendStatus(500)
	}

	// Create a new POST request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return c.SendStatus(500)

	}

	// Set the content type to application/json
	req.Header.Set("Content-Type", "application/json")

	// Perform the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return c.SendStatus(500)

	}
	defer resp.Body.Close()

	// Check the response status
	if resp.StatusCode == http.StatusOK {
		fmt.Println("Request successful")
	} else {
		fmt.Println("Request failed with status:", resp.Status)
	}

	logger.Info("Finished process to send webhook notification")

	return c.Status(200).JSON(fiber.Map{
		"message": "Webhook sent successfully!",
	})
}

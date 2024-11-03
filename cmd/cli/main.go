package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"github.com/tiago123456789/my-own-faas-golang-platform/internal/cli/commands"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var rootCommand = &cobra.Command{
		Use:   "app",
		Short: "CLI - My own faas",
		Long:  "The cli to help to deploy lambda function on my own platform",
	}

	rootCommand.AddCommand(commands.NewDeployCommand().Get())
	rootCommand.Execute()

}

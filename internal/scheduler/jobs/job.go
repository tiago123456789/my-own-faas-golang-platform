package jobs

import (
	"fmt"

	"github.com/tiago123456789/my-own-faas-golang-platform/internal/scheduler/services"
	"github.com/tiago123456789/my-own-faas-golang-platform/pkg/queue"
)

func Init(functionService *services.FunctionService) {

	consumer := queue.NewConsumer(
		"lambda_executions_triggered",
		func(message map[string]interface{}) error {
			fmt.Println(fmt.Sprintf("starting function %s now", message["name"]))
			err := functionService.Run(
				fmt.Sprintf("%s", message["name"]),
				fmt.Sprintf("%s", message["cpu"]),
				fmt.Sprintf("%s", message["memory"]),
			)

			if err != nil {
				fmt.Println(fmt.Sprintf("%v", err))
				return err
			}

			fmt.Println(fmt.Sprintf("Finished function %s", message["name"]))
			return nil
		},
	)

	fmt.Println("Loaded the consumer")
	consumer.Start()
}

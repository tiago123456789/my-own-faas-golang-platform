package jobs

import (
	"fmt"

	"github.com/tiago123456789/my-own-faas-golang-platform/internal/proxy/services"
	"github.com/tiago123456789/my-own-faas-golang-platform/pkg/queue"
)

func Init(functionExecutor services.FunctionExecutorService) {

	consumer := queue.NewConsumer(
		"delete_function_with_expire",
		func(message map[string]interface{}) error {
			fmt.Println("passed on here")
			err := functionExecutor.Stop(message["name"].(string))
			if err != nil {
				return err
			}
			return nil
		})

	go consumer.Start()

}

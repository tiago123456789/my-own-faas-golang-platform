package job

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/tiago123456789/my-own-faas-golang-platform/internal/builder/types"
	"github.com/tiago123456789/my-own-faas-golang-platform/pkg/queue"
)

func Init() {

	publisher := queue.NewPublisher("build_docker_image_process")
	consumer := queue.NewConsumer("builder_docker_image")

	consumer.Consumer(func(message map[string]interface{}) error {
		fmt.Printf("%v", message)

		err := publisher.Publish(types.BuildProgress{
			ID:     fmt.Sprintf("%v", message["id"]),
			Status: "IN_PROGRESS",
		}, 1)

		if err != nil {
			fmt.Printf("Error: %v", err)
			return err
		}
		blueprint := strings.Split(fmt.Sprintf("%s", message["runtime"]), ":")[0]
		blueprintPath := fmt.Sprintf("./internal/builder/blueprint/%s", blueprint)
		fmt.Println(blueprint, blueprintPath)

		srcFolder := fmt.Sprintf("%s", message["path"])
		destFolder := fmt.Sprintf("%s/code.zip", blueprintPath)
		fmt.Println(fmt.Sprintf("cp -rf %s %s", srcFolder, destFolder))
		cpCmd := exec.Command("bash", "-c", fmt.Sprintf("cp -rf %s %s", srcFolder, destFolder))
		err = cpCmd.Run()

		if err != nil {
			fmt.Printf("Error: %v", err)
		}

		commandToBuild := fmt.Sprintf(
			"cd %s && docker build --build-arg MODULE_PATH=%s  --build-arg VERSION_TAG=%s -t tiagorosadacosta123456/lambda-%s .",
			blueprintPath,
			message["moduleName"],
			message["runtime"],
			message["name"],
		)

		cmd := exec.Command("bash", "-c", commandToBuild)
		stdout, _ := cmd.StdoutPipe()
		cmd.Start()
		scanner := bufio.NewScanner(stdout)
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() {
			m := scanner.Text()

			fmt.Println(m)
		}

		err = os.Remove(destFolder)
		if err != nil {
			fmt.Printf("Error: %v", err)
		}

		if err != nil {
			fmt.Println(fmt.Sprintf("Error: %v", err))
			return err
		}

		err = publisher.Publish(types.BuildProgress{
			ID:     fmt.Sprintf("%v", message["id"]),
			Status: "DONE",
		}, 1)
		if err != nil {
			fmt.Printf("Error: %v", err)
			return err
		}
		fmt.Println("Finished the process to build docker image")
		return nil
	})

	consumer.Start()

}

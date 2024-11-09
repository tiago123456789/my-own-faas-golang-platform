package job

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/tiago123456789/my-own-faas-golang-platform/internal/builder/repositories"
	"github.com/tiago123456789/my-own-faas-golang-platform/pkg/queue"
)

func Init(respository *repositories.FunctionRepository) {

	consumer := queue.NewConsumer(
		"builder_docker_image",
		func(message map[string]interface{}) error {
			respository.UpdateProcess(
				message["id"],
				"IN_PROGRESS",
			)

			blueprint := strings.Split(fmt.Sprintf("%s", message["runtime"]), ":")[0]
			blueprintPath := fmt.Sprintf("./internal/builder/blueprint/%s", blueprint)

			srcFolder := fmt.Sprintf("%s", message["path"])
			destFolder := fmt.Sprintf("%s/code.zip", blueprintPath)
			cpCmd := exec.Command("/bin/sh", "-c", fmt.Sprintf("cp -rf %s %s", srcFolder, destFolder))
			err := cpCmd.Run()

			if err != nil {
				fmt.Printf("Error: %v", err)
				respository.UpdateProcess(
					message["id"],
					"FAILED",
				)
				return nil
			}

			verstionTag := fmt.Sprintf("%s", message["runtime"])
			if message["trigger"] == "cron" {
				verstionTag = fmt.Sprintf(
					"%s:%s",
					strings.Split(fmt.Sprintf("%s", message["runtime"]), "-")[0],
					strings.Split(fmt.Sprintf("%s", message["runtime"]), ":")[1],
				)
			}

			commandToBuild := fmt.Sprintf(
				"cd %s && docker build --build-arg MODULE_PATH=%s  --build-arg VERSION_TAG=%s -t tiagorosadacosta123456/lambda-%s .",
				blueprintPath,
				message["moduleName"],
				verstionTag,
				message["name"],
			)

			cmd := exec.Command("/bin/sh", "-c", commandToBuild)
			stdout, _ := cmd.StdoutPipe()

			err = cmd.Start()
			if err != nil {
				fmt.Printf("Error: %v", err)
				respository.UpdateProcess(
					message["id"],
					"FAILED",
				)
				return nil
			}

			scanner := bufio.NewScanner(stdout)
			scanner.Split(bufio.ScanLines)
			for scanner.Scan() {
				m := scanner.Text()
				fmt.Println(m)
			}

			if err := cmd.Wait(); err != nil {
				fmt.Printf("Error: %v", err)
				respository.UpdateProcess(
					message["id"],
					"FAILED",
				)
				return nil
			}

			err = os.Remove(destFolder)
			if err != nil {
				fmt.Printf("Error: %v", err)
				respository.UpdateProcess(
					message["id"],
					"FAILED",
				)
				return nil
			}

			respository.UpdateProcess(
				message["id"],
				"DONE",
			)
			fmt.Println("Finished the process to build docker image")
			return nil
		})

	consumer.Start()

}

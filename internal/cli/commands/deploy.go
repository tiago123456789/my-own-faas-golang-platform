package commands

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"os"
	"time"

	"os/exec"

	"github.com/spf13/cobra"
	httpclient "github.com/tiago123456789/my-own-faas-golang-platform/internal/cli/http-client"
	"github.com/tiago123456789/my-own-faas-golang-platform/internal/cli/types"
	"gopkg.in/yaml.v3"

	modfile "golang.org/x/mod/modfile"
)

var faasUrl string
var runtimesAllowed map[string]bool

type DeployCommand struct {
	httpClient httpclient.HttpClient
}

func init() {
	runtimesAllowed = map[string]bool{
		"golang:1.20":      true,
		"golang:1.19":      true,
		"golang:1.23":      true,
		"golang-cron:1.23": true,
		"golang-cron:1.20": true,
		"golang-cron:1.19": true,
	}

}

func NewDeployCommand() *DeployCommand {
	httpclient := httpclient.New()
	return &DeployCommand{
		httpClient: *httpclient,
	}
}

func (cP *DeployCommand) monitorBuildProgress(functionId string, functionName string) {
	fmt.Println("Starting to build....")

	var responseDataProgress map[string]interface{}
	ticker := time.NewTicker(2 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				urlToSeeProgress := fmt.Sprintf("%s/functions/%s", faasUrl, functionId)
				cP.httpClient.Get(urlToSeeProgress, &responseDataProgress)
				if responseDataProgress["buildProgress"] == "IN_PROGRESS" {
					fmt.Println("Building....")
				}

				if responseDataProgress["buildProgress"] == "DONE" {
					fmt.Printf("The function %s is ready to use!!!\n", functionName)
					ticker.Stop()
					close(quit)
				}
			}
		}
	}()

	<-quit
}

func (cP *DeployCommand) Get() *cobra.Command {
	var path string
	faasUrl = os.Getenv("FAAS_URL")

	var cmd = &cobra.Command{
		Use:   "deploy",
		Short: "Deploy the lambda function",
		Run: func(cmd *cobra.Command, args []string) {
			if path == "" {
				fmt.Println("You need provide the lambda path ")
				return
			}

			_, err := os.Stat(fmt.Sprintf("%s/config.yml", path))
			if os.IsNotExist(err) {
				fmt.Println("Error: to deploy the lambda function you need to create the file config.yml has the configuration of lambda")
				return
			}

			_, err = exec.Command("bash", "-c", fmt.Sprintf("cd %s && zip -r ./code.zip .", path)).Output()
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			yamlFile, err := ioutil.ReadFile(path + "/config.yml")
			if err != nil {
				log.Fatalf("Failed to read YAML file: %v", err)
			}

			var config types.Config
			err = yaml.Unmarshal(yamlFile, &config)
			if err != nil {
				log.Fatalf("The file config.yml has sometime wrong, please the indentention")
			}

			if config.Function.Trigger["http"] == nil && config.Function.Trigger["cron"] == nil {
				log.Fatalf("The triggers supported are: http and cron")
			}

			goModBytes, err := ioutil.ReadFile(fmt.Sprintf("%s/go.mod", path))
			if err != nil {
				fmt.Printf("Error: %v", err)
			}

			modName := modfile.ModulePath(goModBytes)

			if runtimesAllowed[config.Runtime] == false {
				log.Fatalf("The runtine specificed is not valid")
			}

			var data bytes.Buffer
			writer := multipart.NewWriter(&data)

			var trigger string
			if config.Function.Trigger["http"] != nil {
				trigger = "http"
			} else {
				trigger = "cron"
			}

			if len(config.Function.Trigger["cron"]) > 0 {
				err = writer.WriteField("interval", config.Function.Trigger["cron"]["interval"])
				if err != nil {
					fmt.Println("Error writing field:", err)
					return
				}
			}

			err = writer.WriteField("trigger", trigger)
			if err != nil {
				fmt.Println("Error writing field:", err)
				return
			}

			err = writer.WriteField("cpu", config.Cpu)
			if err != nil {
				fmt.Println("Error writing field:", err)
				return
			}

			err = writer.WriteField("memory", config.Memory)
			if err != nil {
				fmt.Println("Error writing field:", err)
				return
			}

			err = writer.WriteField("moduleName", modName)
			if err != nil {
				fmt.Println("Error writing field:", err)
				return
			}

			err = writer.WriteField("runtime", config.Runtime)
			if err != nil {
				fmt.Println("Error writing field:", err)
				return
			}

			err = writer.WriteField("name", config.Name)
			if err != nil {
				fmt.Println("Error writing field:", err)
				return
			}

			file, err := os.Open(fmt.Sprintf("%s/code.zip", path))
			if err != nil {
				fmt.Println("Error opening file:", err)
				return
			}
			defer file.Close()

			part, err := writer.CreateFormFile("file", file.Name())
			if err != nil {
				fmt.Println("Error creating form file:", err)
				return
			}

			_, err = io.Copy(part, file)
			if err != nil {
				fmt.Println("Error copying file:", err)
				return
			}

			err = writer.Close()
			if err != nil {
				fmt.Println("Error closing writer:", err)
				return
			}

			var responseData map[string]interface{}
			err = cP.httpClient.PostMultiPart(
				faasUrl+"/functions", data, writer,
				&responseData,
			)
			if err != nil {
				fmt.Println("Error when tried to deploy the lambda function:", err)
				return
			}

			cP.monitorBuildProgress(
				responseData["id"].(string),
				config.Name,
			)
			return
		},
	}

	cmd.Flags().StringVarP(&path, "path", "p", "", "The lambda function path")
	return cmd
}

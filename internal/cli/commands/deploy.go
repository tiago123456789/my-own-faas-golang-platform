package commands

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"os"

	"os/exec"

	"github.com/spf13/cobra"
	httpclient "github.com/tiago123456789/my-own-faas-golang-platform/internal/cli/http-client"
	"github.com/tiago123456789/my-own-faas-golang-platform/internal/cli/types"
	"gopkg.in/yaml.v3"

	modfile "golang.org/x/mod/modfile"
)

type DeployCommand struct {
	httpClient httpclient.HttpClient
}

func NewDeployCommand() *DeployCommand {
	httpclient := httpclient.New()
	return &DeployCommand{
		httpClient: *httpclient,
	}
}

func (cP *DeployCommand) Get() *cobra.Command {
	var path string

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

			if config.Function.Trigger["http"] == nil {
				log.Fatalf("The only trigger supported is http")
			}

			goModBytes, err := ioutil.ReadFile(fmt.Sprintf("%s/go.mod", path))
			if err != nil {
				fmt.Printf("Error: %v", err)
			}

			modName := modfile.ModulePath(goModBytes)

			runtimesAllowed := []string{
				"golang:1.20", "golang:1.19", "golang:1.23",
			}

			for _, runtime := range runtimesAllowed {
				if runtime == config.Runtime {
					var data bytes.Buffer
					writer := multipart.NewWriter(&data)

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

					err = cP.httpClient.PostMultiPart(
						"http://localhost:3000/functions", data, writer,
					)
					if err != nil {
						fmt.Println("Error when tried to deploy the lambda function:", err)
						return
					}

					fmt.Println("Execute with success")
					return
				}

			}

			log.Fatalf("The runtine specificed is not valid")
		},
	}

	cmd.Flags().StringVarP(&path, "path", "p", "", "The lambda function path")
	return cmd
}

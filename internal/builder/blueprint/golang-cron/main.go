package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	lambda "MODULE_NAME"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/yaml.v3"
)

type Function struct {
	Trigger map[string]map[string]string `yaml:"trigger"`
}

type Config struct {
	Function Function          `yaml:"function"`
	Name     string            `yaml:"name"`
	Envs     map[string]string `yaml:"envs"`
}

// HTTPWriter is a custom writer that sends logs to an HTTP endpoint.
type HTTPWriter struct {
	URL string
}

// Write sends the log entry to the configured HTTP endpoint.
func (h *HTTPWriter) Write(p []byte) (n int, err error) {
	req, err := http.NewRequest("POST", h.URL, bytes.NewBuffer(p))
	if err != nil {
		return 0, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, err
	}

	return len(p), nil
}

func main() {

	httpWriter := &HTTPWriter{
		URL: "http://host.docker.internal:5050/",
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder, // e.g. INFO, WARN
		EncodeTime:     zapcore.ISO8601TimeEncoder,  // e.g. 2021-01-01T00:00:00Z
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(httpWriter),
		zapcore.InfoLevel,
	)

	pathCode := os.Getenv("PATH_CODE")
	yamlFile, err := ioutil.ReadFile(pathCode + "/config.yml")
	if err != nil {
		log.Fatalf("Failed to read YAML file: %v", err)
	}

	var config Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		log.Fatalf("Failed to unmarshal YAML: %v", err)
	}

	if len(config.Envs) > 0 {
		for key, value := range config.Envs {
			err := os.Setenv(key, value)
			if err != nil {
				fmt.Printf("Error setting environment variable: %v\n", err)
				return
			}
		}
	}

	logger := zap.New(core).With(
		zap.String("service", config.Name),
	)

	logger.Info("Starting execution the function " + config.Name)
	start := time.Now()
	lambda.Handler(logger)
	timeTook := time.Since(start)
	logger.Info("Finished execution the function " + config.Name)
	logger.Info("The function " + config.Name + " spent " + timeTook.String() + " seconds")
	time.AfterFunc(2*time.Second, func() {
		os.Exit(0)
	})

}

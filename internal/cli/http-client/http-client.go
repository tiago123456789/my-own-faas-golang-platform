package httpclient

import (
	"bytes"
	"encoding/json"
	"log"
	"mime/multipart"
	"net/http"
)

type HttpClient struct {
}

func New() *HttpClient {
	return &HttpClient{}
}

func (h *HttpClient) Get(url string, data interface{}) error {
	response, err := http.Get(url)
	if err != nil {
		log.Fatalf("Error making GET request: %v", err)
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		log.Fatalf("Error: received status code %d", response.StatusCode)
	}

	if err := json.NewDecoder(response.Body).Decode(&data); err != nil {
		log.Fatalf("Error decoding JSON: %v", err)
		return err
	}

	return nil
}

func (h *HttpClient) PostMultiPart(
	url string,
	data bytes.Buffer,
	multipartWriter *multipart.Writer,
	responseData interface{},
) error {
	req, err := http.NewRequest("POST", url, &data)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", multipartWriter.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&responseData); err != nil {
		log.Fatalf("Error decoding JSON: %v", err)
		return err
	}

	return nil
}

package httpclient

import (
	"bytes"
	"mime/multipart"
	"net/http"
)

type HttpClient struct {
}

func New() *HttpClient {
	return &HttpClient{}
}

func (h *HttpClient) PostMultiPart(
	url string,
	data bytes.Buffer,
	multipartWriter *multipart.Writer,
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

	return nil
}

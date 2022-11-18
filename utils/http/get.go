package http

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type (
	Client interface {
		Do(req *http.Request) (*http.Response, error)
	}
)

func Get[T any](client Client, url string) (T, error) {
	var empty T
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return empty, fmt.Errorf("failed to create http request: %v", err)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return empty, fmt.Errorf("failed to send http request: %v", err)
	}
	defer resp.Body.Close()

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return empty, fmt.Errorf("failed to read response content: %v", err)
	}

	var result T
	if jErr := json.Unmarshal(content, &result); jErr != nil {
		return empty, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	return result, nil
}

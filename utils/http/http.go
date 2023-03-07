package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type (
	Client interface {
		Do(req *http.Request) (*http.Response, error)
	}
)

var EmptyQuery map[string]string

func Get[T any](client Client, url string, query map[string]string) (T, error) {
	var empty T
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return empty, fmt.Errorf("failed to create http request: %v", err)
	}
	req.Header.Set("Accept", "application/json")

	if len(query) > 0 {
		q := req.URL.Query()
		for key, value := range query {
			q.Add(key, value)
		}
		req.URL.RawQuery = q.Encode()
	}

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
		return empty, fmt.Errorf("failed to unmarshal response: %v", jErr)
	}

	return result, nil
}

func Post[T any](client Client, url string, body any) (T, error) {
	var (
		empty T
		req   *http.Request
		err   error
	)
	if body != nil {
		marshaledBody, marshallErr := json.Marshal(body)
		if marshallErr != nil {
			return empty, fmt.Errorf("failed to marshall body: %v", marshallErr)
		}
		req, err = http.NewRequest("POST", url, bytes.NewBuffer(marshaledBody))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, err = http.NewRequest("POST", url, nil)
	}
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
		return empty, fmt.Errorf("failed to unmarshal response: %v", jErr)
	}

	return result, nil
}

func IsValidUrl(str string) bool {
	parsedUrl, err := url.Parse(str)
	return err == nil && parsedUrl.Scheme != "" && parsedUrl.Host != ""
}

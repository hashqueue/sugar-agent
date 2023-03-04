package utils

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

// HTTPClient 是自定义的HTTP客户端结构体
type HTTPClient struct {
	BaseURL string
}

// Get 发送GET请求
func (c *HTTPClient) Get(path string, headers map[string]string) ([]byte, error) {
	url := c.BaseURL + path
	req, err := http.NewRequest("GET", url, nil)
	FailOnError(err, "Failed to create request")

	// 添加请求头
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := http.DefaultClient.Do(req)
	FailOnError(err, "Failed to send request")
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		FailOnError(err, "Failed to close response body")
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	FailOnError(err, "Failed to read response body")
	return body, nil
}

// Post 发送POST请求
func (c *HTTPClient) Post(path string, headers map[string]string, body interface{}) ([]byte, error) {
	url := c.BaseURL + path

	jsonData, err := json.Marshal(body)
	FailOnError(err, "Failed to marshal body")

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	FailOnError(err, "Failed to create request")

	// 添加请求头
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	FailOnError(err, "Failed to send request")
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		FailOnError(err, "Failed to close response body")
	}(resp.Body)

	respBody, err := io.ReadAll(resp.Body)
	FailOnError(err, "Failed to read response body")
	return respBody, nil
}

// Put 发送PUT请求
func (c *HTTPClient) Put(path string, headers map[string]string, body interface{}) ([]byte, error) {
	url := c.BaseURL + path

	jsonData, err := json.Marshal(body)
	FailOnError(err, "Failed to marshal body")

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	FailOnError(err, "Failed to create request")

	// 添加请求头
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	FailOnError(err, "Failed to send request")
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		FailOnError(err, "Failed to close response body")
	}(resp.Body)

	respBody, err := io.ReadAll(resp.Body)
	FailOnError(err, "Failed to read response body")

	return respBody, nil
}

// Patch 发送PATCH请求
func (c *HTTPClient) Patch(path string, headers map[string]string, body interface{}) ([]byte, error) {
	url := c.BaseURL + path

	jsonData, err := json.Marshal(body)
	FailOnError(err, "Failed to marshal body")

	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonData))
	FailOnError(err, "Failed to create request")

	// 添加请求头
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	FailOnError(err, "Failed to send request")
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		FailOnError(err, "Failed to close response body")
	}(resp.Body)

	respBody, err := io.ReadAll(resp.Body)
	FailOnError(err, "Failed to read response body")
	return respBody, nil
}

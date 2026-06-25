package helpers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/yourorg/go-mvc-app/src/config"
)

type APIClient struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
	headers    map[string]string
}

type APIResponse struct {
	StatusCode int
	Body       []byte
	Headers    http.Header
}

func NewAPIClient() *APIClient {
	cfg := config.App.API
	return &APIClient{
		baseURL: cfg.BaseURL,
		apiKey:  cfg.Key,
		httpClient: &http.Client{
			Timeout: time.Duration(cfg.Timeout) * time.Second,
		},
		headers: map[string]string{
			"Content-Type": "application/json",
			"Accept":       "application/json",
			"X-API-Key":    cfg.Key,
		},
	}
}

// WithHeader adds or overrides a default header for the client.
func (c *APIClient) WithHeader(key, value string) *APIClient {
	c.headers[key] = value
	return c
}

// Get performs a GET request to path and decodes the JSON response into dest.
func (c *APIClient) Get(ctx context.Context, path string, dest interface{}) (*APIResponse, error) {
	return c.do(ctx, http.MethodGet, path, nil, dest)
}

// Post performs a POST request with a JSON body and decodes the response into dest.
func (c *APIClient) Post(ctx context.Context, path string, body interface{}, dest interface{}) (*APIResponse, error) {
	return c.do(ctx, http.MethodPost, path, body, dest)
}

// Put performs a PUT request with a JSON body and decodes the response into dest.
func (c *APIClient) Put(ctx context.Context, path string, body interface{}, dest interface{}) (*APIResponse, error) {
	return c.do(ctx, http.MethodPut, path, body, dest)
}

// Patch performs a PATCH request with a JSON body and decodes the response into dest.
func (c *APIClient) Patch(ctx context.Context, path string, body interface{}, dest interface{}) (*APIResponse, error) {
	return c.do(ctx, http.MethodPatch, path, body, dest)
}

// Delete performs a DELETE request and decodes the response into dest.
func (c *APIClient) Delete(ctx context.Context, path string, dest interface{}) (*APIResponse, error) {
	return c.do(ctx, http.MethodDelete, path, nil, dest)
}

func (c *APIClient) do(ctx context.Context, method, path string, body interface{}, dest interface{}) (*APIResponse, error) {
	var reqBody io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("api_client: marshal body: %w", err)
		}
		reqBody = bytes.NewBuffer(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, reqBody)
	if err != nil {
		return nil, fmt.Errorf("api_client: build request: %w", err)
	}

	for k, v := range c.headers {
		req.Header.Set(k, v)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("api_client: execute request: %w", err)
	}
	defer resp.Body.Close()

	rawBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("api_client: read response: %w", err)
	}

	apiResp := &APIResponse{
		StatusCode: resp.StatusCode,
		Body:       rawBody,
		Headers:    resp.Header,
	}

	if resp.StatusCode >= 400 {
		return apiResp, fmt.Errorf("api_client: upstream error %d: %s", resp.StatusCode, string(rawBody))
	}

	if dest != nil && len(rawBody) > 0 {
		if err := json.Unmarshal(rawBody, dest); err != nil {
			return apiResp, fmt.Errorf("api_client: decode response: %w", err)
		}
	}

	return apiResp, nil
}

package hightouch

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"terraform-provider-hightouch/pkg/helper"
	"time"
)

// Client is a client for the Hightouch API.
type Client struct {
	apiKey     string
	httpClient *http.Client
	baseURL    string
}

// APIError represents an error response from the Hightouch API.
type APIError struct {
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

func (e APIError) Error() string {
	return fmt.Sprintf("Hightouch API Error: %s", e.Message)
}

// makeRequest is a helper function to create, send, and handle API requests.
func (c *Client) makeRequest(method, path string, body interface{}) ([]byte, error) {
	var reqBody io.Reader
	fmt.Print(body)
	if body != nil {
		jsonBytes, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonBytes)

		// Peek into the request body for debugging purposes
		peek, err := helper.PrettyPrintJson(string(jsonBytes))
		fmt.Println(peek)
	}

	url := fmt.Sprintf("%s%s", c.baseURL, path)
	fmt.Printf("%s: %s\n", method, url)
	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set necessary headers
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "hightouch_client-go-client/1.0")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Printf("Warning: failed to close response body: %v\n", err)
		}
	}(resp.Body)

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	peek, err := helper.PrettyPrintJson(string(respBody))
	fmt.Println(peek)

	// Check for non-successful status codes
	fmt.Printf("Response: %s\n", string(respBody))
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var apiErr APIError
		if err := json.Unmarshal(respBody, &apiErr); err != nil {
			// If we can't parse the error, return a generic one
			return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(respBody))
		}
		return nil, apiErr
	}

	return respBody, nil
}

// NewClient creates a new Hightouch API client.
// It requires an API key, which can be generated from your Hightouch workspace settings.
func NewClient(apiKey string, apiBaseUrl string) *Client {
	return &Client{
		apiKey:     apiKey,
		httpClient: &http.Client{Timeout: 15 * time.Second},
		baseURL:    apiBaseUrl,
	}
}

package hightouch

import (
	"encoding/json"
	"fmt"
	"time"
)

type HightouchSource struct {
	ID            *int                   `json:"id"`
	Name          string                 `json:"name"`
	Slug          string                 `json:"slug"`
	WorkspaceID   int                    `json:"workspaceId"`
	CreatedAt     time.Time              `json:"createdAt"`
	UpdatedAt     time.Time              `json:"updatedAt"`
	Type          string                 `json:"type"`
	Configuration map[string]interface{} `json:"configuration"`
}

// GetSnowflakeSource GetHightouchSource retrieves a specific source by its ID.
func (c *Client) GetSnowflakeSource(
	sourceID int,
) (*HightouchSource, error) {

	var source HightouchSource

	respBody, err := c.makeRequest(
		"GET",
		fmt.Sprintf("/sources/%d", sourceID),
		nil,
	)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(respBody, &source); err != nil {
		return nil, fmt.Errorf("failed to unmarshal GetHightouchSource response: %w", err)
	}

	return &source, nil
}

func (c *Client) CreateHightouchSource(
	name string,
	slug string,
	sourceType string,
	configuration map[string]interface{},
) (*HightouchSource, error) {
	requestBody := struct {
		Name          string                 `json:"name"`
		Slug          string                 `json:"slug"`
		Type          string                 `json:"type"`
		Configuration map[string]interface{} `json:"configuration"`
	}{
		Name:          name,
		Slug:          slug,
		Type:          sourceType,
		Configuration: configuration,
	}

	var source HightouchSource
	respBody, err := c.makeRequest(
		"POST",
		"/sources",
		requestBody,
	)

	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(respBody, &source); err != nil {
		return nil, fmt.Errorf("failed to unmarshal CreateHightouchSource response: %w", err)
	}

	return &source, nil
}

// UpdateHightouchSource updates a specific source.
// The `updates` map can contain any of the mutable source fields, e.g., "name", "configuration".
func (c *Client) UpdateHightouchSource(
	sourceID int,
	name string,
	configuration map[string]interface{},
) (*HightouchSource, error) {
	requestBody := struct {
		Name          string                 `json:"name"`
		Configuration map[string]interface{} `json:"configuration"`
	}{
		Name:          name,
		Configuration: configuration,
	}

	var source HightouchSource
	respBody, err := c.makeRequest(
		"PATCH",
		fmt.Sprintf("/sources/%d", sourceID),
		requestBody,
	)

	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(respBody, &source); err != nil {
		return nil, fmt.Errorf("failed to unmarshal UpdateHightouchSource response: %w", err)
	}

	return &source, nil
}

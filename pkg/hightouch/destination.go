package hightouch

import (
	"encoding/json"
	"fmt"
	"time"
)

type HightouchDestination struct {
	ID            *int                   `json:"id"`
	Name          string                 `json:"name"`
	Slug          string                 `json:"slug"`
	WorkspaceID   int                    `json:"workspaceId"`
	CreatedAt     time.Time              `json:"createdAt"`
	UpdatedAt     time.Time              `json:"updatedAt"`
	Type          string                 `json:"type"`
	Syncs         []int                  `json:"syncs"`
	Configuration map[string]interface{} `json:"configuration"`
}

// GetHightouchDestination retrieves a specific destination by its ID.
func (c *Client) GetHightouchDestination(
	destinationID int,
) (*HightouchDestination, error) {

	var destination HightouchDestination

	respBody, err := c.makeRequest(
		"GET",
		fmt.Sprintf("/destinations/%d", destinationID),
		nil,
	)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(respBody, &destination); err != nil {
		return nil, fmt.Errorf("failed to unmarshal GetHightouchDestination response: %w", err)
	}

	return &destination, nil
}

// CreateHightouchDestination creates a new destination in Hightouch.
func (c *Client) CreateHightouchDestination(
	name string,
	slug string,
	destinationType string,
	configuration map[string]interface{},
) (*HightouchDestination, error) {
	requestBody := struct {
		Name          string                 `json:"name"`
		Slug          string                 `json:"slug"`
		Type          string                 `json:"type"`
		Configuration map[string]interface{} `json:"configuration"`
	}{
		Name:          name,
		Slug:          slug,
		Type:          destinationType,
		Configuration: configuration,
	}

	var destination HightouchDestination
	respBody, err := c.makeRequest(
		"POST",
		"/destinations",
		requestBody,
	)

	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(respBody, &destination); err != nil {
		return nil, fmt.Errorf("failed to unmarshal CreateHightouchDestination response: %w", err)
	}

	return &destination, nil
}

// UpdateHightouchDestination updates a specific destination.
// The name and configuration parameters can be updated.
func (c *Client) UpdateHightouchDestination(
	destinationID int,
	name string,
	configuration map[string]interface{},
) (*HightouchDestination, error) {
	requestBody := struct {
		Name          string                 `json:"name"`
		Configuration map[string]interface{} `json:"configuration"`
	}{
		Name:          name,
		Configuration: configuration,
	}

	var destination HightouchDestination
	respBody, err := c.makeRequest(
		"PATCH",
		fmt.Sprintf("/destinations/%d", destinationID),
		requestBody,
	)

	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(respBody, &destination); err != nil {
		return nil, fmt.Errorf("failed to unmarshal UpdateHightouchDestination response: %w", err)
	}

	return &destination, nil
}

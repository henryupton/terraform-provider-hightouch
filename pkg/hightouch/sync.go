package hightouch

import (
	"encoding/json"
	"fmt"
	"time"
)

type HightouchSync struct {
	ID            *int                   `json:"id"`
	Name          string                 `json:"name"`
	Slug          string                 `json:"slug"`
	WorkspaceID   int                    `json:"workspaceId"`
	CreatedAt     time.Time              `json:"createdAt"`
	UpdatedAt     time.Time              `json:"updatedAt"`
	DestinationID int                    `json:"destinationId"`
	ModelID       int                    `json:"modelId"`
	Configuration map[string]interface{} `json:"configuration"`
	Schedule      map[string]interface{} `json:"schedule"`
	Status        string                 `json:"status"`
	Disabled      bool                   `json:"disabled"`
}

// GetHightouchSync retrieves a specific sync by its ID.
func (c *Client) GetHightouchSync(
	syncID int,
) (*HightouchSync, error) {

	var sync HightouchSync

	respBody, err := c.makeRequest(
		"GET",
		fmt.Sprintf("/syncs/%d", syncID),
		nil,
	)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(respBody, &sync); err != nil {
		return nil, fmt.Errorf("failed to unmarshal GetHightouchSync response: %w", err)
	}

	return &sync, nil
}

// CreateHightouchSync creates a new sync in Hightouch.
func (c *Client) CreateHightouchSync(
	name string,
	slug string,
	sourceID int,
	destinationID int,
	modelID int,
	configuration map[string]interface{},
	schedule map[string]interface{},
) (*HightouchSync, error) {
	requestBody := struct {
		Name          string                 `json:"name"`
		Slug          string                 `json:"slug"`
		SourceID      int                    `json:"sourceId"`
		DestinationID int                    `json:"destinationId"`
		ModelID       int                    `json:"modelId"`
		Configuration map[string]interface{} `json:"configuration"`
		Schedule      map[string]interface{} `json:"schedule"`
	}{
		Name:          name,
		Slug:          slug,
		SourceID:      sourceID,
		DestinationID: destinationID,
		ModelID:       modelID,
		Configuration: configuration,
		Schedule:      schedule,
	}

	var sync HightouchSync
	respBody, err := c.makeRequest(
		"POST",
		"/syncs",
		requestBody,
	)

	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(respBody, &sync); err != nil {
		return nil, fmt.Errorf("failed to unmarshal CreateHightouchSync response: %w", err)
	}

	return &sync, nil
}

// UpdateHightouchSync updates a specific sync.
// The name, configuration, schedule, and disabled status can be updated.
func (c *Client) UpdateHightouchSync(
	syncID int,
	name string,
	configuration map[string]interface{},
	schedule map[string]interface{},
	disabled bool,
) (*HightouchSync, error) {
	requestBody := struct {
		Name          string                 `json:"name"`
		Configuration map[string]interface{} `json:"configuration"`
		Schedule      map[string]interface{} `json:"schedule"`
		Disabled      bool                   `json:"disabled"`
	}{
		Name:          name,
		Configuration: configuration,
		Schedule:      schedule,
		Disabled:      disabled,
	}

	var sync HightouchSync
	respBody, err := c.makeRequest(
		"PATCH",
		fmt.Sprintf("/syncs/%d", syncID),
		requestBody,
	)

	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(respBody, &sync); err != nil {
		return nil, fmt.Errorf("failed to unmarshal UpdateHightouchSync response: %w", err)
	}

	return &sync, nil
}

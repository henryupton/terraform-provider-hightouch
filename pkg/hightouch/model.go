package hightouch

import (
	"encoding/json"
	"fmt"
	"time"
)

type HightouchModel struct {
	ID          *int      `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	WorkspaceID int       `json:"workspaceId"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	SourceID    int       `json:"sourceId"`
	SQL         string    `json:"sql"`
	DBTable     string    `json:"dbtTable"`
	QueryType   string    `json:"queryType"`
	PrimaryKey  string    `json:"primaryKey"`
	IsSchema    bool      `json:"isSchema"`
}

// GetHightouchModel retrieves a specific model by its ID.
func (c *Client) GetHightouchModel(
	modelID int,
) (*HightouchModel, error) {

	var model HightouchModel

	respBody, err := c.makeRequest(
		"GET",
		fmt.Sprintf("/models/%d", modelID),
		nil,
	)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(respBody, &model); err != nil {
		return nil, fmt.Errorf("failed to unmarshal GetHightouchModel response: %w", err)
	}

	return &model, nil
}

// CreateHightouchModel creates a new model in Hightouch.
func (c *Client) CreateHightouchModel(
	name string,
	slug string,
	sourceID int,
	sql string,
	queryType string,
	primaryKey string,
) (*HightouchModel, error) {
	requestBody := struct {
		Name       string `json:"name"`
		Slug       string `json:"slug"`
		SourceID   int    `json:"sourceId"`
		SQL        string `json:"sql"`
		QueryType  string `json:"queryType"`
		PrimaryKey string `json:"primaryKey"`
	}{
		Name:       name,
		Slug:       slug,
		SourceID:   sourceID,
		SQL:        sql,
		QueryType:  queryType,
		PrimaryKey: primaryKey,
	}

	var model HightouchModel
	respBody, err := c.makeRequest(
		"POST",
		"/models",
		requestBody,
	)

	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(respBody, &model); err != nil {
		return nil, fmt.Errorf("failed to unmarshal CreateHightouchModel response: %w", err)
	}

	return &model, nil
}

// UpdateHightouchModel updates a specific model.
// The name, SQL, primary key, description, and tags can be updated.
func (c *Client) UpdateHightouchModel(
	modelID int,
	name string,
	sql string,
	primaryKey string,
) (*HightouchModel, error) {
	requestBody := struct {
		Name       string `json:"name"`
		SQL        string `json:"sql"`
		PrimaryKey string `json:"primaryKey"`
	}{
		Name:       name,
		SQL:        sql,
		PrimaryKey: primaryKey,
	}

	var model HightouchModel
	respBody, err := c.makeRequest(
		"PATCH",
		fmt.Sprintf("/models/%d", modelID),
		requestBody,
	)

	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(respBody, &model); err != nil {
		return nil, fmt.Errorf("failed to unmarshal UpdateHightouchModel response: %w", err)
	}

	return &model, nil
}

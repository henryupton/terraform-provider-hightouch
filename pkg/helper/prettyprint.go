package helper

import "encoding/json"

// PrettyPrintJson takes a JSON string and returns a formatted version.
func PrettyPrintJson(rawJson string) (string, error) {
	var data interface{}
	// Unmarshal the raw JSON into a generic interface
	if err := json.Unmarshal([]byte(rawJson), &data); err != nil {
		return "", err
	}

	// Marshal the data back into an indented format
	prettyBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", err
	}

	return string(prettyBytes), nil
}

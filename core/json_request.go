package core

import (
	"encoding/json"
	"fmt"
)

// JSONRequest is the base request.
type JSONRequest struct {
	Path    string
	Request BaseRequest
}

const jsonContentType = "application/json"

// Validate the request.
func (r JSONRequest) Validate() error {
	return r.Request.Validate()
}

// EndpointPath returns with the path for the endpoint.
func (r JSONRequest) EndpointPath() string {
	return r.Path
}

// ToBody returns with the JSON []byte representation of the request.
func (r JSONRequest) ToBody(token string) ([]byte, string, error) {
	requestBody, err := json.Marshal(r.Request)
	if err != nil {
		return requestBody, jsonContentType, fmt.Errorf("unable to encode request: %w", err)
	}

	var repack map[string]interface{}

	if err := json.Unmarshal(requestBody, &repack); err != nil {
		return requestBody, jsonContentType, fmt.Errorf("failed to decode body: %w", err)
	}

	if token != "" {
		repack["i"] = token
	}

	content, err := json.Marshal(repack)
	if err != nil {
		return content, jsonContentType, fmt.Errorf("failed to encode body: %w", err)
	}

	return content, jsonContentType, nil
}

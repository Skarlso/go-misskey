package antennas

import "github.com/yitsushi/go-misskey/core"

// DeleteRequest represents a request to delete an antenna.
type DeleteRequest struct {
	AntennaID string `json:"antennaId"`
}

// Delete is the endpoint to delete an existing antenna.
func (s *Service) Delete(antennaID string) error {
	request := &DeleteRequest{
		AntennaID: antennaID,
	}

	var response core.DummyResponse
	err := s.Call(
		&core.BaseRequest{Request: request, Path: "/antennas/delete"},
		&response,
	)

	return err
}

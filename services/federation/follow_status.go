package federation

import (
	"github.com/yitsushi/go-misskey/core"
	"github.com/yitsushi/go-misskey/models"
)

const (
	followingPath = "/federation/following"
	followersPath = "/federation/followers"
)

// FollowStatusList contains a list of status informations.
type FollowStatusList struct {
	Items []models.FollowStatus
}

// FollowStatusRequest contains request information obtain the status of followers or followees.
type FollowStatusRequest struct {
	Host    string `json:"host"`
	SinceID string `json:"sinceId"`
	UntilID string `json:"untilId"`
	Limit   int    `json:"limit"`
}

// Validate the request.
func (r *FollowStatusRequest) Validate() error {
	if r.Host == "" {
		return core.RequestValidationError{
			Request: r,
			Message: core.UndefinedRequiredField,
			Field:   "Host",
		}
	}
	return nil
}

// Followers lists all followers.
func (s *Service) Followers(request FollowStatusRequest) (FollowStatusList, error) {
	return s.call(request, followersPath)
}

// Following lists all followings.
func (s *Service) Following(request FollowStatusRequest) (FollowStatusList, error) {
	return s.call(request, followingPath)
}

// call will make the call to the service with the given path and request.
func (s *Service) call(request FollowStatusRequest, path string) (FollowStatusList, error) {
	var response FollowStatusList

	err := s.Call(
		&core.JSONRequest{Request: &request, Path: path},
		&response,
	)

	return response, err
}
package youtube

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/theartofdevel/youtube-search-service/pkg/client/youtube"
	"github.com/theartofdevel/youtube-search-service/pkg/logging"
	"net/http"
)

type Service struct {
	Client youtube.Client
	Logger *logging.Logger
}

func NewService(client youtube.Client, logger *logging.Logger) *Service {
	return &Service{Client: client, Logger: logger}
}

func (s *Service) UpdateAccessToken(ctx context.Context) error {
	response, err := s.Client.UpdateAccessToken(ctx)

	var uat youtube.UpdateAccessTokenResponse
	if err = json.NewDecoder(response.Body).Decode(&uat); err != nil {
		return fmt.Errorf("failed to decode body due to err: %v", err)
	}

	s.Client.SetAccessToken(uat.AccessToken)
	s.Client.SetRefreshToken(uat.RefreshToken)

	return nil
}

func (s *Service) FindTrackByName(ctx context.Context, trackName string) (string, error) {
	response, err := s.Client.SearchTrack(ctx, trackName)
	if err != nil {
		return "", err
	}

	if response.StatusCode == http.StatusUnauthorized || response.StatusCode == http.StatusForbidden {
		err := s.UpdateAccessToken(ctx)
		if err != nil {
			return "", err
		}

		response, err = s.Client.SearchTrack(ctx, trackName)
		if err != nil {
			return "", err
		}
	}

	var responseData map[string]interface{}
	if err = json.NewDecoder(response.Body).Decode(&responseData); err != nil {
		return "", err
	}

	if a, ok := responseData["items"].([]interface{}); ok {
		b := a[0].(map[string]interface{})["id"].(map[string]interface{})["videoId"].(string)
		return fmt.Sprintf("https://youtube.com/watch?v=%s", b), nil
	} else {
		return "", fmt.Errorf("youtube request failed due to error %v", responseData)
	}
}

package youtube

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/theartofdevel/youtube-search-service/pkg/client/youtube"
	"github.com/theartofdevel/youtube-search-service/pkg/logging"
)

type service struct {
	client youtube.Client
	logger *logging.Logger
}

func NewService(client youtube.Client, logger *logging.Logger) Service {
	return &service{client: client, logger: logger}
}

type Service interface {
	FindTrackByName(ctx context.Context, trackName string) (string, error)
}

func (s *service) FindTrackByName(ctx context.Context, trackName string) (string, error) {
	response, err := s.client.SearchTrack(ctx, trackName)
	if err != nil {
		return "", err
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

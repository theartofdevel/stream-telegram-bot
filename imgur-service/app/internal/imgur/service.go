package imgur

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/theartofdevel/imgur-service/pkg/client/imgur"
	"github.com/theartofdevel/imgur-service/pkg/logging"
)

type imgurService struct {
	client imgur.Client
	logger *logging.Logger
}

func NewImgurService(client imgur.Client, logger *logging.Logger) Service {
	return &imgurService{client: client, logger: logger}
}

type Service interface {
	ShareImage(ctx context.Context, image []byte) (string, error)
}

func (i *imgurService) ShareImage(ctx context.Context, image []byte) (string, error) {
	response, err := i.client.UploadImage(ctx, image)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	var responseData map[string]interface{}
	if err = json.NewDecoder(response.Body).Decode(&responseData); err != nil {
		return "", err
	}

	if response.StatusCode != 200 {
		i.logger.Error(responseData)
		return "", fmt.Errorf("failed to upload image")
	}

	return responseData["data"].(map[string]interface{})["link"].(string), nil
}

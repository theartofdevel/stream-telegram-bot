package youtube

import (
	"encoding/json"
	"fmt"
	"github.com/theartofdevel/telegram_bot/internal/events"
	"github.com/theartofdevel/telegram_bot/internal/events/model"
	"github.com/theartofdevel/telegram_bot/pkg/logging"
)

type yt struct {
	logger *logging.Logger
}

func NewYouTubeProcessEventStrategy(logger *logging.Logger) events.ProcessEventStrategy {
	return &yt{
		logger: logger,
	}
}

func (p *yt) Process(eventBody []byte) (response model.ProcessedEvent, err error) {
	event := SearchTrackResponse{}
	if err = json.Unmarshal(eventBody, &event); err != nil {
		return response, fmt.Errorf("failed to unmarshal event due to error %v", err)
	}
	
	if event.Meta.Error != nil {
		response.Err = fmt.Errorf(*event.Meta.Error)
	}
	response.RequestID = event.Meta.RequestID
	response.Message = event.Data.URL

	return response, nil
}

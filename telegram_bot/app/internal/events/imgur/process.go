package imgur

import (
	"encoding/json"
	"fmt"
	"github.com/theartofdevel/telegram_bot/internal/events"
	"github.com/theartofdevel/telegram_bot/internal/events/model"
	"github.com/theartofdevel/telegram_bot/pkg/logging"
)

type imgur struct {
	logger *logging.Logger
}

func NewImgurProcessEventStrategy(logger *logging.Logger) events.ProcessEventStrategy {
	return &imgur{
		logger: logger,
	}
}

func (p *imgur) Process(eventBody []byte) (response model.ProcessedEvent, err error) {
	event := UploadImageResponse{}
	if err = json.Unmarshal(eventBody, &event); err != nil {
		return response, fmt.Errorf("failed to unmarshal event due to error %v", err)
	}
	var eventErr error
	if event.Meta.Error != nil {
		eventErr = fmt.Errorf(*event.Meta.Error)
	}
	return model.ProcessedEvent{
		RequestID: event.Meta.RequestID,
		Message:   event.Data.URL,
		Err:       eventErr,
	}, nil
}

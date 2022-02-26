package youtube

import (
	"github.com/theartofdevel/telegram_bot/internal/events/model"
)

type SearchTrackRequest struct {
	RequestID string `json:"request_id"`
	Name      string `json:"name"`
}

type SearchTrackResponse struct {
	Meta model.ResponseMeta `json:"meta"`
	Data ResponseData       `json:"data"`
}
type ResponseData struct {
	URL string `json:"name"`
}

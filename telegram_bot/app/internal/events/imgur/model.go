package imgur

import (
	"github.com/theartofdevel/telegram_bot/internal/events/model"
)

type UploadImageRequest struct {
	RequestID string `json:"request_id"`
	Photo     []byte `json:"photo"`
}

type UploadImageResponse struct {
	Meta model.ResponseMeta `json:"meta"`
	Data ResponseData       `json:"data"`
}
type ResponseData struct {
	URL string `json:"name"`
}

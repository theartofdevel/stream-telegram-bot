package request

type UploadImage struct {
	RequestID string `json:"request_id"`
	Photo     []byte `json:"photo"`
}

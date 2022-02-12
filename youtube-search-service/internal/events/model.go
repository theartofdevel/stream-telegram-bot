package events

type SearchTrack struct {
	RequestID string `json:"request_id"`
	Name      string `json:"name"`
}

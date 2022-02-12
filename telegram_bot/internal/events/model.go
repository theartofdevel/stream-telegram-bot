package events

type SearchTrackRequest struct {
	RequestID string `json:"request_id"`
	Name      string `json:"name"`
}

type SearchTrackResponse struct {
	RequestID string `json:"request_id,omitempty"`
	Name      string `json:"name,omitempty"`
	Success   string `json:"success,omitempty"`
	Error     string `json:"err,omitempty"`
}

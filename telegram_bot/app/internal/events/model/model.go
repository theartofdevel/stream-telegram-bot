package model

import "fmt"

type ResponseMessage struct {
	Meta ResponseMeta `json:"meta"`
	Data interface{}  `json:"data"`
}

type ResponseMeta struct {
	RequestID string  `json:"request_id"`
	Error     *string `json:"err,omitempty"`
}

func (m *ResponseMeta) String() string {
	return fmt.Sprintf("RequestID:%s, Error:%s", m.RequestID, m.Error)
}

type ProcessedEvent struct {
	RequestID string
	Message   string
	Err       error
}

func (m *ProcessedEvent) String() string {
	return fmt.Sprintf(
		"RequestID:%s. Message:%s. Err:%s", m.RequestID, m.Message, m.Err,
	)
}

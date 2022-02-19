package events

import (
	"context"
	"encoding/json"
	"github.com/theartofdevel/youtube-search-service/internal/events/model/request"
	"github.com/theartofdevel/youtube-search-service/internal/events/model/response"
	"github.com/theartofdevel/youtube-search-service/internal/youtube"
	"github.com/theartofdevel/youtube-search-service/pkg/client/mq"
	"github.com/theartofdevel/youtube-search-service/pkg/logging"
)

type worker struct {
	id            int
	client        mq.Consumer
	producer      mq.Producer
	responseQueue string
	messages      <-chan mq.Message
	logger        *logging.Logger
	service       youtube.Service
}

type Worker interface {
	Process()
}

func NewWorker(id int, client mq.Consumer, responseQueue string, producer mq.Producer, messages <-chan mq.Message, logger *logging.Logger, service youtube.Service) Worker {
	return &worker{id: id, client: client, responseQueue: responseQueue, messages: messages, producer: producer, logger: logger, service: service}
}

func (w *worker) Process() {
	for msg := range w.messages {
		event := request.SearchTrack{}
		if err := json.Unmarshal(msg.Body, &event); err != nil {
			w.logger.Errorf("[worker #%d]: failed to unmarchal event due to error %v", w.id, err)
			w.logger.Debugf("[worker #%d]: body: %s", w.id, msg.Body)

			w.reject(msg)
			continue
		}

		respData := response.SearchTrack{
			Meta: response.Meta{
				RequestID: event.RequestID,
			},
			Data: response.Data{},
		}
		var errorStr string
		url, err := w.service.FindTrackByName(context.TODO(), event.Name)
		if err != nil {
			errorStr = err.Error()
			respData.Meta.Error = &errorStr
		} else {
			respData.Data.URL = url
		}

		w.sendResponse(respData)

		w.ack(msg)
	}
}

func (w *worker) sendResponse(d interface{}) {
	b, err := json.Marshal(d)
	if err != nil {
		w.logger.Errorf("[worker #%d]: failed to response due to error %v", w.id, err)
		return
	}
	err = w.producer.Publish(w.responseQueue, b)
	if err != nil {
		w.logger.Errorf("[worker #%d]: failed to response due to error %v", w.id, err)
	}
}

func (w *worker) reject(msg mq.Message) {
	if err := w.client.Reject(msg.ID, false); err != nil {
		w.logger.Errorf("[worker #%d]: failed to reject due to error %v", w.id, err)
	}
}

func (w *worker) ack(msg mq.Message) {
	if err := w.client.Ack(msg.ID, false); err != nil {
		w.logger.Errorf("[worker #%d]: failed to ack due to error %v", w.id, err)
	}
}

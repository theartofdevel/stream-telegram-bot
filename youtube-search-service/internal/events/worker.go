package events

import (
	"context"
	"encoding/json"
	"github.com/theartofdevel/youtube-search-service/internal/youtube"
	"github.com/theartofdevel/youtube-search-service/pkg/client/mq"
	"github.com/theartofdevel/youtube-search-service/pkg/logging"
	"strconv"
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
		event := SearchTrack{}
		if err := json.Unmarshal(msg.Body, &event); err != nil {
			w.logger.Errorf("[worker #%d]: failed to unmarchal event due to error %v", w.id, err)
			w.logger.Debugf("[worker #%d]: body: %s", w.id, msg.Body)

			w.reject(msg)
			continue
		}

		respData := map[string]string{
			"request_id": event.RequestID,
		}
		name, err := w.service.FindTrackByName(context.TODO(), event.Name)
		if err != nil {
			respData["err"] = err.Error()
		} else {
			respData["name"] = name
		}

		respData["success"] = strconv.FormatBool(err == nil)
		w.sendResponse(respData)

		w.ack(msg)
	}
}

func (w *worker) sendResponse(d map[string]string) {
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

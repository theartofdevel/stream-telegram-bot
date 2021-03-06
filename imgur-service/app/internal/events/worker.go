package events

import (
	"context"
	"encoding/json"
	"github.com/theartofdevel/imgur-service/internal/events/model/request"
	"github.com/theartofdevel/imgur-service/internal/events/model/response"
	"github.com/theartofdevel/imgur-service/internal/imgur"
	"github.com/theartofdevel/imgur-service/pkg/client/mq"
	"github.com/theartofdevel/imgur-service/pkg/logging"
	"sync"
)

type Worker struct {
	id            int
	client        mq.Consumer
	producer      mq.Producer
	responseQueue string
	messages      <-chan mq.Message
	logger        *logging.Logger
	service       imgur.Service
	wg            *sync.WaitGroup
}

func NewWorker(id int, client mq.Consumer, responseQueue string, producer mq.Producer, messages <-chan mq.Message,
	logger *logging.Logger, service imgur.Service, wg *sync.WaitGroup) Worker {
	return Worker{id: id, client: client, responseQueue: responseQueue, messages: messages, producer: producer,
		logger: logger, service: service, wg: wg}
}

func (w *Worker) Process() {
	defer w.wg.Done()
	for msg := range w.messages {
		event := request.UploadImage{}
		if err := json.Unmarshal(msg.Body, &event); err != nil {
			w.logger.Errorf("[Worker #%d]: failed to unmarshal event due to error %v", w.id, err)
			w.logger.Debugf("[Worker #%d]: body: %s", w.id, msg.Body)

			w.reject(msg)
			continue
		}

		respData := response.UploadImage{
			Meta: response.Meta{
				RequestID: event.RequestID,
			},
			Data: response.Data{},
		}
		var errorStr string
		url, err := w.service.ShareImage(context.TODO(), event.Photo)
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

func (w *Worker) sendResponse(d interface{}) {
	b, err := json.Marshal(d)
	if err != nil {
		w.logger.Errorf("[Worker #%d]: failed to response due to error %v", w.id, err)
		return
	}
	err = w.producer.Publish(w.responseQueue, b)
	if err != nil {
		w.logger.Errorf("[Worker #%d]: failed to response due to error %v", w.id, err)
	}
}

func (w *Worker) reject(msg mq.Message) {
	if err := w.client.Reject(msg.ID, false); err != nil {
		w.logger.Errorf("[Worker #%d]: failed to reject due to error %v", w.id, err)
	}
}

func (w *Worker) ack(msg mq.Message) {
	if err := w.client.Ack(msg.ID, false); err != nil {
		w.logger.Errorf("[Worker #%d]: failed to ack due to error %v", w.id, err)
	}
}

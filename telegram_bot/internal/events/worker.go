package events

import (
	"encoding/json"
	"github.com/theartofdevel/telegram_bot/pkg/client/mq"
	"github.com/theartofdevel/telegram_bot/pkg/logging"
	tele "gopkg.in/telebot.v3"
	"strconv"
)

type worker struct {
	id            int
	client        mq.Consumer
	producer      mq.Producer
	responseQueue string
	messages      <-chan mq.Message
	logger        *logging.Logger
	bot           *tele.Bot
}

type Worker interface {
	Process()
}

func NewWorker(id int, client mq.Consumer, bot *tele.Bot, producer mq.Producer, messages <-chan mq.Message, logger *logging.Logger) Worker {
	return &worker{id: id, client: client, bot: bot, messages: messages, producer: producer, logger: logger}
}

func (w *worker) Process() {
	for msg := range w.messages {
		event := SearchTrackResponse{}
		if err := json.Unmarshal(msg.Body, &event); err != nil {
			w.logger.Errorf("[worker #%d]: failed to unmarchal event due to error %v", w.id, err)
			w.logger.Debugf("[worker #%d]: body: %s", w.id, msg.Body)

			w.reject(msg)
			continue
		}

		i, _ := strconv.ParseInt(event.RequestID, 10, 64)
		id, err := w.bot.ChatByID(i)
		if err != nil {
			w.logger.Errorf("[worker #%d]: failed to get chat by id due to error %v", w.id, err)
		}

		message := "Запрос не обработан, произошла ошибка"
		if event.Success == "true" {
			message = event.Name
		}

		_, err = w.bot.Send(id, message)
		if err != nil {
			w.logger.Errorf("[worker #%d]: failed to get send due to error %v", w.id, err)
		}
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

package events

import (
	"github.com/theartofdevel/telegram_bot/internal/service/bot"
	"github.com/theartofdevel/telegram_bot/pkg/client/mq"
	"github.com/theartofdevel/telegram_bot/pkg/logging"
)

type worker struct {
	id              int
	client          mq.Consumer
	producer        mq.Producer
	responseQueue   string
	messages        <-chan mq.Message
	logger          *logging.Logger
	processStrategy ProcessEventStrategy
	botService      bot.Service
}

type Worker interface {
	Process()
}

func NewWorker(id int, client mq.Consumer, processStrategy ProcessEventStrategy, botService bot.Service, producer mq.Producer, messages <-chan mq.Message, logger *logging.Logger) Worker {
	return &worker{id: id, client: client, processStrategy: processStrategy, botService: botService, messages: messages, producer: producer, logger: logger}
}

func (w *worker) Process() {
	for msg := range w.messages {
		processedEvent, err := w.processStrategy.Process(msg.Body)
		if err != nil {
			w.logger.Errorf("[worker #%d]: failed to processedEvent event due to error %v", w.id, err)
			w.logger.Debugf("[worker #%d]: body: %s", w.id, msg.Body)
			w.reject(msg)
			return
		}

		err = w.botService.SendMessage(processedEvent)
		if err != nil {
			w.logger.Errorf("[worker #%d]: failed to sent message due to error %v", w.id, err)
			w.logger.Debugf("[worker #%d]: body: %s", w.id, msg.Body)
			w.reject(msg)
			return
		}

		w.ack(msg)
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

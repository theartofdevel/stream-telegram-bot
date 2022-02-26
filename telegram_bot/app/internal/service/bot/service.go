package bot

import (
	"fmt"
	"github.com/theartofdevel/telegram_bot/internal/events/model"
	"github.com/theartofdevel/telegram_bot/pkg/logging"
	tele "gopkg.in/telebot.v3"
	"strconv"
)

type Service struct {
	Bot    *tele.Bot
	Logger *logging.Logger
}

func (bs *Service) SendMessage(data model.ProcessedEvent) error {
	i, _ := strconv.ParseInt(data.RequestID, 10, 64)
	id, err := bs.Bot.ChatByID(i)
	if err != nil {
		bs.Logger.Tracef("Bot Send ResponseMessage. ProcessedEvent: %s", data)
		return fmt.Errorf("failed to get chat by id due to error %v", err)
	}

	message := data.Message
	if data.Err != nil {
		message = fmt.Sprintf("Запрос не обработан, произошла ошибка (%s)", data.Err)
	}

	_, err = bs.Bot.Send(id, message)
	if err != nil {
		bs.Logger.Tracef("ChatID: %d, Data: %s", id.ID, data)
		return fmt.Errorf("failed to get send due to error %v", err)
	}

	return nil
}

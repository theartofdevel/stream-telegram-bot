package rabbitmq

import (
	"fmt"
	"github.com/streadway/amqp"
	"github.com/theartofdevel/youtube-search-service/pkg/client/mq"
)

type ProducerConfig struct {
	BaseConfig
}

type rabbitMQProducer struct {
	*rabbitMQBase
}

func NewRabbitMQProducer(cfg ProducerConfig) (mq.Producer, error) {
	producer := &rabbitMQProducer{
		rabbitMQBase: &rabbitMQBase{
			done: make(chan bool),
		},
	}

	addr := fmt.Sprintf("amqp://%s:%s@%s:%s/", cfg.Username, cfg.Password, cfg.Host, cfg.Port)
	err := producer.connect(addr)
	if err != nil {
		return nil, err
	}

	go producer.handleReconnect(addr)

	return producer, nil
}

func (r *rabbitMQProducer) Publish(target string, body []byte) error {
	if !r.Connected() {
		return errNotConnected
	}

	err := r.ch.Publish(
		"",
		target,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         body,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message due %v", err)
	}

	return nil
}

func (r *rabbitMQProducer) Close() error {
	if err := r.close(); err != nil {
		return err
	}

	return nil
}

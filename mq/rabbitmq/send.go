package rabbitmq

import (
	"fmt"
	"github.com/streadway/amqp"
	"github.com/vgmdj/utils/logger"
)

func (mq *rabbit) SendToQue(exchange, key string, body []byte) (err error) {
	err = mq.ch.ExchangeDeclare(exchange, "direct", true, false, false, false, nil)
	if err != nil {
		logger.Error(fmt.Sprintf("%s: %s\n", "Failed to declare a exchange", err))
		return
	}

	err = mq.ch.Publish(
		exchange, // exchange
		key,      // routing key
		false,    // mandatory
		false,    // immediate
		amqp.Publishing{
			ContentType:  "text/plain",
			Body:         body,
			DeliveryMode: amqp.Persistent,
		})
	logger.Info(fmt.Sprintf(" [x] Sent %s", string(body)))
	if err != nil {
		logger.Error(fmt.Sprintf("%s: %s\n", "Failed to publish a message", err))
		return
	}

	return
}

func (mq *rabbit) SendToDLQue(exchange, key, queue string, body []byte, args amqp.Table) (err error) {
	err = mq.ch.ExchangeDeclare(exchange, "direct", true, false, false, false, nil)
	if err != nil {
		logger.Error(fmt.Sprintf("%s: %s\n", "Failed to declare a exchange", err))
		return
	}

	_, err = mq.ch.QueueDeclare(queue, true, false, false, false, args)
	if err != nil {
		logger.Error(fmt.Sprintf("%s: %s\n", "Failed to declare a queue", err))
		return
	}

	err = mq.ch.QueueBind(queue, key, exchange, false, nil)
	if err != nil {
		logger.Error(fmt.Sprintf("%s: %s\n", "Failed to bind a exchange", err))
		return
	}

	err = mq.ch.Publish(
		exchange, // exchange
		key,      // routing key
		false,    // mandatory
		false,    // immediate
		amqp.Publishing{
			ContentType:  "text/plain",
			Body:         body,
			DeliveryMode: amqp.Persistent,
		})
	logger.Info(fmt.Sprintf(" [x] Sent %s", body))
	if err != nil {
		logger.Error(fmt.Sprintf("%s: %s\n", "Failed to publish a message", err))
		return
	}

	return
}

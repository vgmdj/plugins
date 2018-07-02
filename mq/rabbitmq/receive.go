package rabbitmq

import (
	"fmt"
	"github.com/streadway/amqp"
	"github.com/vgmdj/utils/logger"
)

//GetQueue assume the queue is already exist
func (mq *rabbit) GetQueue(name string, args amqp.Table) amqp.Queue {
	queue, err := mq.ch.QueueDeclarePassive(name, true, false,
		false, false, args)
	if err != nil {
		logger.Error(err.Error())
		return queue
	}

	return queue
}

//ReceiveFromMQ
func (mq *rabbit) ReceiveFromMQ(exchange, key, queue string, args amqp.Table) (msgs <-chan amqp.Delivery, err error) {
	if mq.ch == nil {
		err = fmt.Errorf("Failed to connect to RabbitMQ")
		return
	}

	err = mq.ch.ExchangeDeclare(exchange, "direct", true, false, false, false, nil)
	if err != nil {
		logger.Error("%s: %s\n", "Failed to declare a exchange", err)
		return
	}

	_, err = mq.ch.QueueDeclare(queue, true, false, false, false, args)
	if err != nil {
		logger.Error("%s: %s\n", "Failed to declare a queue", err)
		return
	}

	err = mq.ch.QueueBind(queue, key, exchange, false, nil)
	if err != nil {
		logger.Error("%s: %s\n", "Failed to bind queue to exchange", err)
		return
	}

	msgs, err = mq.ch.Consume(
		queue, // queue
		queue, // consumer
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		logger.Error("%s: %s\n", "Failed to register a consumer", err)
		return
	}

	return
}

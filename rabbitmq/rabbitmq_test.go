package rabbitmq

import (
	"github.com/vgmdj/utils/logger"
	"log"
	"testing"
	"time"
)

func TestNewRabbit(t *testing.T) {
	rabbit, _ := NewRabbit("127.0.0.1:5672", "/",
		"user", "pwd")

	//send to mq
	rabbit.SendToQue("exchange", "key", []byte("OK"))

	//receive from mq
	msgs, _ := rabbit.ReceiveFromMQ("exchange", "key", "queue", nil)
	for msg := range msgs {
		log.Println("Receive a message from mq: ", string(msg.Body))

		//do some thing

	}

}

func TestPing(t *testing.T) {
	logger.SetLogFuncCall(true)

	rabbit, _ := NewRabbit("10.11.22.101:31099", "/",
		"wangrui", "vgmdj")

	//send to mq
	rabbit.SendToQue("exchange", "key", []byte("OK"))

	//check health
	go func() {
		for {
			time.Sleep(time.Second)

			queue, err := rabbit.GetQueue("queue", nil)
			if err == nil {
				logger.Info(queue)
				continue
			}

			err = rabbit.Reconnect()
			if err != nil {
				logger.Error(err.Error())
				continue
			}

			go receive(rabbit)

		}
	}()

	rabbit.SetQos(1, 0, true)
	go receive(rabbit)

	time.Sleep(time.Hour)

}

func receive(rabbit *Rabbit) {
	//receive from mq
	msgs, _ := rabbit.ReceiveFromMQ("exchange", "key", "queue", nil)
	for msg := range msgs {
		logger.Info("Receive a message from mq: ", string(msg.Body))

		//do some thing

		time.Sleep(time.Second * 5)

		logger.Info("ack")
		if err := msg.Ack(false); err != nil {
			logger.Error(err.Error())
			return
		}

	}
}

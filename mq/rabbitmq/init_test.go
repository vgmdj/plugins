package rabbitmq

import (
	"log"
	"testing"
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

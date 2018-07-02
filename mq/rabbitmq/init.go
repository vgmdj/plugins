package rabbitmq

import (
	"fmt"
	"github.com/streadway/amqp"
	"github.com/vgmdj/utils/logger"
	"sync"
)

//map[dialURL]*amqp.Connection
var connection sync.Map

//Rabbit
type rabbit struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

//NewRabbit
func NewRabbit(server, vhost, userName, password string) (*rabbit, error) {

	if vhost == "" {
		vhost = "/"
	} else if vhost[0] != '/' {
		vhost = "/" + vhost
	}

	var (
		connValue interface{}
		ok        bool
		conn      *amqp.Connection
		ch        *amqp.Channel
		err       error
	)

	dialURL := fmt.Sprintf("amqp://%s:%s@%s%s", userName, password, server, vhost)
	if connValue, ok = connection.Load(dialURL); !ok {
		return connect(dialURL)
	}

	conn, ok = connValue.(*amqp.Connection)
	if !ok {
		logger.Error("sync map error")
		return connect(dialURL)
	}

	ch, err = conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("%s: %s\n", "Failed to open a channel", err)
	}

	return &rabbit{conn: conn, ch: ch}, nil

}

//connect
func connect(dialURL string) (*rabbit, error) {
	conn, err := amqp.Dial(dialURL)
	if err != nil {
		return nil, fmt.Errorf("%s: %s\n", "Failed to connect to RabbitMQ", err)
	}

	connection.Store(dialURL, conn)

	var ch *amqp.Channel
	ch, err = conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("%s: %s\n", "Failed to open a channel", err)
	}

	return &rabbit{
		conn: conn,
		ch:   ch,
	}, nil
}

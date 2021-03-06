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
type Rabbit struct {
	conn    *amqp.Connection
	ch      *amqp.Channel
	dialURL string
}

//NewRabbit
func NewRabbit(server, vhost, userName, password string) (*Rabbit, error) {

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
		connection.Delete(dialURL)
		return nil, fmt.Errorf("%s: %s\n", "Failed to open a channel", err)
	}

	return &Rabbit{conn: conn, ch: ch, dialURL: dialURL}, nil

}

//Reconnect
func (mq *Rabbit) Reconnect() (err error) {

	conn, err := amqp.Dial(mq.dialURL)
	if err != nil {
		return fmt.Errorf("%s: %s\n", "Failed to connect to RabbitMQ", err)
	}

	connection.Store(mq.dialURL, conn)

	var ch *amqp.Channel
	ch, err = conn.Channel()
	if err != nil {
		return fmt.Errorf("%s: %s\n", "Failed to open a channel", err)
	}

	mq.conn = conn
	mq.ch = ch

	return
}

//SetQos
func (mq *Rabbit) SetQos(count, size int, global bool) (err error) {
	return mq.ch.Qos(count, size, global)
}

//connect
func connect(dialURL string) (*Rabbit, error) {
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

	return &Rabbit{
		conn:    conn,
		ch:      ch,
		dialURL: dialURL,
	}, nil
}

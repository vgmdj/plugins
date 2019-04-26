package redis

import (
	"bytes"
	"errors"
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
)

//Message represents a message notification.
type Message struct {
	//Kind is "subscribe", "unsubscribe", "psubscribe" or "punsubscribe"
	Kind string

	//The originating channel.
	Channel string

	//The message data.
	Data []byte
}

//Pong represents a pubsub pong notification.
type Pong struct {
	Data string
}

//PubSub publish and subscribe
type PubSub struct {
	conn redis.Conn
}

//NewSession return redips PubSub
func (c *Client) NewPubSub() *PubSub {
	return &PubSub{
		conn: c.pool.Get(),
	}
}

//Publish publish messages to specified channel and return subscribe
func (ps *PubSub) Publish(channel string, arg interface{}) (int, error) {
	return redis.Int(ps.conn.Do("PUBLISH", channel, arg))
}

//Subscribe subscribes the connection to the specified channels.
func (ps *PubSub) Subscribe(channel ...interface{}) error {
	ps.conn.Send("SUBSCRIBE", channel...)
	return ps.conn.Flush()
}

//PSubscribe subscribes the connection to the given patterns.
func (ps PubSub) PSubscribe(channel ...interface{}) error {
	ps.conn.Send("PSUBSCRIBE", channel...)
	return ps.conn.Flush()
}

//Unsubscribe unsubscribes the connection from the given channels, or from all
//of them if none is given.
func (ps PubSub) Unsubscribe(channel ...interface{}) error {
	ps.conn.Send("UNSUBSCRIBE", channel...)
	return ps.conn.Flush()
}

//PUnsubscribe unsubscribes the connection from the given patterns, or from all
//of them if none is given.
func (ps PubSub) PUnsubscribe(channel ...interface{}) error {
	ps.conn.Send("PUNSUBSCRIBE", channel...)
	return ps.conn.Flush()
}

//Ping sends a PING to the server with the specified data.
//
//The connection must be subscribed to at least one channel or pattern when
//calling this method.
func (ps PubSub) Ping(data string) error {
	ps.conn.Send("PING", data)
	return ps.conn.Flush()
}

//Close closes the connection.
func (ps PubSub) Close() error {
	return ps.conn.Close()
}

//Receive returns a pushed message as a Subscription, Message, Pong or error.
//The return value is intended to be used directly in a type switch as
//illustrated in the PubSubConn example.
func (ps PubSub) Receive() (Message, error) {
	return ps.receiveInternal(ps.conn.Receive())
}

//ReceiveWithTimeout is like Receive, but it allows the application to
//override the connection's default timeout.
func (ps PubSub) ReceiveWithTimeout(timeout time.Duration) (Message, error) {
	return ps.receiveInternal(redis.ReceiveWithTimeout(ps.conn, timeout))
}

func (ps PubSub) receiveInternal(replyArg interface{}, errArg error) (m Message, err error) {
	reply, err := redis.Values(replyArg, errArg)
	if err != nil {
		return
	}

	reply, err = redis.Scan(reply, &m.Kind)
	if err != nil {
		return
	}

	switch m.Kind {
	case "message":
		_, err = redis.Scan(reply, &m.Channel, &m.Data)
		return

	case "pmessage":
		_, err = redis.Scan(reply, &m.Kind, &m.Channel, &m.Data)
		return

	case "subscribe", "psubscribe", "unsubscribe", "punsubscribe":
		count := 0
		_, err = redis.Scan(reply, &m.Channel, &count)
		buf := new(bytes.Buffer)
		fmt.Fprint(buf, count)
		m.Data = buf.Bytes()

		return

	case "pong":
		_, err = redis.Scan(reply, &m.Data)
		return

	}
	return m, errors.New("redigo: unknown pubsub notification")
}

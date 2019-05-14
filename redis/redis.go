package redis

import (
	"context"
	"net"
	"sync"
	"time"

	"github.com/gomodule/redigo/redis"
)

var (
	ErrNil = redis.ErrNil

	DefaultConf = &ClientConf{
		Address:        "127.0.0.1:6379",
		ConnectTimeout: time.Second * 30,
	}

	cli     *Client
	cliOnce sync.Once
)

//ClientConf redis client config
type ClientConf struct {
	Address        string
	Password       string
	DB             int
	ConnectTimeout time.Duration
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	KeepAlive      time.Duration

	// Maximum number of idle connections in the pool.
	MaxIdle int

	// Maximum number of connections allocated by the pool at a given time.
	// When zero, there is no limit on the number of connections in the pool.
	MaxActive int

	// Close connections after remaining idle for this duration. If the value
	// is zero, then idle connections are not closed. Applications should set
	// the timeout to a value less than the server's timeout.
	IdleTimeout time.Duration

	// If Wait is true and the pool is at the MaxActive limit, then Get() waits
	// for a connection to be returned to the pool before returning.
	Wait bool

	// Close connections older than this duration. If the value is zero, then
	// the pool does not close connections based on age.
	MaxConnLifetime time.Duration
}

//init initialize
func (conf *ClientConf) init() {
	// As per the IANA draft spec, the host defaults to localhost and
	// the port defaults to 6379.
	host, port, err := net.SplitHostPort(conf.Address)
	if err != nil {
		// assume port is missing
		host = conf.Address
		port = "6379"
	}
	if host == "" {
		host = "localhost"
	}
	conf.Address = net.JoinHostPort(host, port)
}

//DialOptions return the dial options
func (conf *ClientConf) DialOptions() []redis.DialOption {
	options := make([]redis.DialOption, 0)
	if conf.DB != 0 {
		options = append(options, redis.DialDatabase(conf.DB))
	}

	if conf.Password != "" {
		options = append(options, redis.DialPassword(conf.Password))
	}

	if conf.ConnectTimeout != 0 {
		options = append(options, redis.DialConnectTimeout(conf.ConnectTimeout))
	}

	if conf.WriteTimeout != 0 {
		options = append(options, redis.DialWriteTimeout(conf.WriteTimeout))
	}

	if conf.ReadTimeout != 0 {
		options = append(options, redis.DialReadTimeout(conf.ReadTimeout))
	}

	if conf.KeepAlive != 0 {
		options = append(options, redis.DialKeepAlive(conf.KeepAlive))
	}

	return options
}

//Client redis client
type Client struct {
	pool *redis.Pool
}

//NewClient before connect will ping , default ping timeout is a minute
func NewClient(conf *ClientConf) *Client {
	conf.init()

	return &Client{
		pool: &redis.Pool{
			// Other pool configuration not shown in this example.
			Dial: func() (redis.Conn, error) {
				c, err := redis.Dial("tcp", conf.Address, conf.DialOptions()[:]...)
				if err != nil {
					return nil, err
				}

				return c, nil
			},

			TestOnBorrow: func(c redis.Conn, t time.Time) error {
				if time.Since(t) < time.Minute {
					return nil
				}
				_, err := c.Do("PING")
				return err
			},

			MaxIdle: conf.MaxIdle,

			MaxActive: conf.MaxActive,

			MaxConnLifetime: conf.MaxConnLifetime,

			Wait: conf.Wait,

			IdleTimeout: conf.IdleTimeout,
		},
	}

}

func UniqueClient(conf *ClientConf) *Client {
	cliOnce.Do(func() {
		cli = NewClient(conf)
	})

	return cli
}

//Ping check the redis status
func (c *Client) Ping() (err error) {
	conn := c.pool.Get()
	defer conn.Close()

	_, err = conn.Do("PING")
	return err
}

//Get return redis conn
func (c *Client) GetConn() redis.Conn {
	return c.pool.Get()
}

//GetCtxConn return redis conn with provide context
func (c *Client) GetCtxConn(ctx context.Context) (redis.Conn, error) {
	return c.pool.GetContext(ctx)
}

//Do sends a command to the server and returns the received reply.
func (c *Client) Do(commandName string, args ...interface{}) (reply interface{}, err error) {
	conn := c.pool.Get()
	defer conn.Close()

	return conn.Do(commandName, args[:]...)
}

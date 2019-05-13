package redis

import (
	"fmt"

	"github.com/gomodule/redigo/redis"
)

var (
	ExistErr = fmt.Errorf("key already exist")
)

//Set store string
func (c *Client) Set(key, value interface{}) (err error) {
	conn := c.pool.Get()
	defer conn.Close()

	_, err = conn.Do("SET", key, value)

	return
}

//SetNX set unique string type value
func (c *Client) SetNX(key string, value interface{}, seconds int) (err error) {
	conn := c.pool.Get()
	defer conn.Close()

	_, err = redis.String(conn.Do("SET", key, value, "EX", seconds, "NX"))
	if err == redis.ErrNil {
		return ExistErr
	}

	return

}

//GetString string get string
func (c *Client) Get(key string) (interface{}, error) {
	conn := c.pool.Get()
	defer conn.Close()

	return conn.Do("GET", key)

}

func (c *Client) GetBytes(key string) ([]byte, error) {
	conn := c.pool.Get()
	defer conn.Close()

	return redis.Bytes(conn.Do("GET", key))
}

//GetString string get string
func (c *Client) GetString(key string) (string, error) {
	conn := c.pool.Get()
	defer conn.Close()

	return redis.String(conn.Do("GET", key))

}

//GetInt string get int
func (c *Client) GetInt(key string) (int, error) {
	conn := c.pool.Get()
	defer conn.Close()

	return redis.Int(conn.Do("GET", key))

}

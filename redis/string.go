package redis

import (
	"github.com/garyburd/redigo/redis"
)

//Set store string
func (c *Client) Set(key, value interface{}) (err error) {
	conn := c.pool.Get()
	defer conn.Close()

	if _, err = conn.Do("SET", key, value); err != nil {
		return
	}

	return
}

//SetNX set unique string type value
func (c *Client) SetNX(key string, value interface{}, seconds int) (err error) {
	conn := c.pool.Get()
	defer conn.Close()

	if _, err = redis.String(conn.Do("SET", key, value, "EX", seconds, "NX")); err != nil {
		return
	}

	return

}

//Expire set expire
func (c *Client) Expire(key string, seconds int) (err error) {
	conn := c.pool.Get()
	defer conn.Close()

	if _, err = conn.Do("EXPIRE", key, seconds); err != nil {
		return
	}

	return
}

//GetString string get string
func (c *Client) Get(key string) (interface{}, error) {
	conn := c.pool.Get()
	defer conn.Close()

	return conn.Do("GET", key)

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

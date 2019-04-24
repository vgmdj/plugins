package redis

import (
	"fmt"

	"github.com/gomodule/redigo/redis"
)

type Type string

const (
	TypeNone   Type = "none"
	TypeString Type = "string"
	TypeList   Type = "list"
	TypeSet    Type = "set"
	TypeZSet   Type = "zset"
	TypeHash   Type = "hash"
	TypeStream Type = "stream"
)

var (
	NotExistErr = fmt.Errorf("key not exist")
)

//Expire set expire, set seconds=0 will delete the key
func (c *Client) Expire(key string, seconds int) (err error) {
	if seconds < 0 {
		seconds = 0
	}

	conn := c.pool.Get()
	defer conn.Close()

	if _, err = conn.Do("EXPIRE", key, seconds); err != nil {
		return
	}

	return
}

//Exists return key existence
func (c *Client) Exists(key string) (bool, error) {
	conn := c.pool.Get()
	defer conn.Close()

	return redis.Bool(conn.Do("EXISTS", key))

}

//DEL delete key , is equal to Expire(key,0)
func (c *Client) DEL(key string) error {
	conn := c.pool.Get()
	defer conn.Close()

	_, err := conn.Do("DEL", key)
	if err != nil {
		return err
	}

	return nil

}

//Type return the type of the key
func (c *Client) Type(key string) (Type, error) {
	conn := c.pool.Get()
	defer conn.Close()

	reply, err := redis.String(conn.Do("TYPE", key))
	return Type(reply), err
}

//Rename rename the key
func (c *Client) Rename(key string, newKeyName string) error {
	conn := c.pool.Get()
	defer conn.Close()

	_, err := conn.Do("RENAME", key, newKeyName)
	if err != nil {
		return err
	}

	return nil

}

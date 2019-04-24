package redis

import (
	"testing"
	"time"
)

func testRedisClient() *Client {
	c := NewClient(&ClientConf{
		Address:        "10.11.22.77",
		Password:       "wangrui",
		DB:             4,
		ConnectTimeout: time.Second * 5,
	})

	_, err := c.Do("PING")
	if err != nil {
		panic(err.Error())
	}

	return c
}

func TestNewClient(t *testing.T) {
	testRedisClient()
}

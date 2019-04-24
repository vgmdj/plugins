package redis

import (
	"github.com/gomodule/redigo/redis"
)

//HSet hash set
func (c *Client) HSet(key, itemKey string, item interface{}) (err error) {
	conn := c.pool.Get()
	defer conn.Close()

	if _, err = conn.Do("HSET", key, itemKey, item); err != nil {
		return
	}

	return
}

//HExists select hash key-item's existence
func (c *Client) HExists(key, itemKey string) (bool, error) {
	conn := c.pool.Get()
	defer conn.Close()

	return redis.Bool(conn.Do("HEXISTS", key, itemKey))
}

//HDel hash delete key-item
func (c *Client) HDel(key, itemKey string) (err error) {
	conn := c.pool.Get()
	defer conn.Close()

	_, err = conn.Do("HDEL", key, itemKey)
	if err != nil {
		return
	}

	return
}

//HGet hash get interface
func (c *Client) HGet(userKey, itemKey string) (interface{}, error) {
	conn := c.pool.Get()
	defer conn.Close()

	return redis.Values(conn.Do("HGET", userKey, itemKey))

}

//HGetBool hash get bool
func (c *Client) HGetBool(userKey, itemKey string) (bool, error) {
	conn := c.pool.Get()
	defer conn.Close()

	return redis.Bool(conn.Do("HGET", userKey, itemKey))

}

//HGetBytes hash get []byte
func (c *Client) HGetBytes(userKey, itemKey string) ([]byte, error) {
	conn := c.pool.Get()
	defer conn.Close()

	return redis.Bytes(conn.Do("HGET", userKey, itemKey))

}

//HGetFloat64 hash get float64
func (c *Client) HGetFloat64(userKey, itemKey string) (float64, error) {
	conn := c.pool.Get()
	defer conn.Close()

	return redis.Float64(conn.Do("HGET", userKey, itemKey))

}

//HGetInt get int
func (c *Client) HGetInt(userKey, itemKey string) (int, error) {
	conn := c.pool.Get()
	defer conn.Close()

	return redis.Int(conn.Do("HGET", userKey, itemKey))

}

//HGetInt64 get int64
func (c *Client) HGetInt64(userKey, itemKey string) (int64, error) {
	conn := c.pool.Get()
	defer conn.Close()

	return redis.Int64(conn.Do("HGET", userKey, itemKey))

}

//HGetInts hash get ints
func (c *Client) HGetInts(userKey, itemKey string) ([]int, error) {
	conn := c.pool.Get()
	defer conn.Close()

	return redis.Ints(conn.Do("HGET", userKey, itemKey))

}

//HGetString hash get string
func (c *Client) HGetString(userKey, itemKey string) (string, error) {
	conn := c.pool.Get()
	defer conn.Close()

	return redis.String(conn.Do("HGET", userKey, itemKey))

}

//HGetStrings hash get strings
func (c *Client) HGetStrings(userKey, itemKey string) ([]string, error) {
	conn := c.pool.Get()
	defer conn.Close()

	return redis.Strings(conn.Do("HGET", userKey, itemKey))

}

//HGetIntMap hash get int map
func (c *Client) HGetIntMap(userKey, itemKey string) (map[string]int, error) {
	conn := c.pool.Get()
	defer conn.Close()

	return redis.IntMap(conn.Do("HGET", userKey, itemKey))

}

//HGetInt64Map hash get int64 map
func (c *Client) HGetInt64Map(userKey, itemKey string) (map[string]int64, error) {
	conn := c.pool.Get()
	defer conn.Close()

	return redis.Int64Map(conn.Do("HGET", userKey, itemKey))

}

//HGetStringMap hash get string map
func (c *Client) HGetStringMap(userKey, itemKey string) (map[string]string, error) {
	conn := c.pool.Get()
	defer conn.Close()

	return redis.StringMap(conn.Do("HGET", userKey, itemKey))

}

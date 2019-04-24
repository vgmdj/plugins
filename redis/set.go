package redis

import (
	"github.com/garyburd/redigo/redis"
	"github.com/vgmdj/utils/logger"
)

//SAdd and key-value and return the count of key set
func (c *Client) SAdd(key string, args ...interface{}) (count int64, err error) {
	conn := c.pool.Get()
	defer conn.Close()

	for _, v := range args {
		if count, err = redis.Int64(conn.Do("SADD", key, v)); err != nil {
			return
		}
	}

	return
}

//SCard return the count of key set
func (c *Client) SCard(key string) (int64, error) {
	conn := c.pool.Get()
	defer conn.Close()

	return redis.Int64(conn.Do("SCARD", key))
}

//SRem remove the member of key set
func (c *Client) SRem(key string, member interface{}) (err error) {
	conn := c.pool.Get()
	defer conn.Close()

	_, err = conn.Do("SREM", key, member)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	return
}

//SMembers return all members of key set
func (c *Client) SMembers(key string) ([]interface{}, error) {
	conn := c.pool.Get()
	defer conn.Close()

	return redis.Values(conn.Do("SMEMBERS", key))

}

//SMembersInts return all int members of key set
func (c *Client) SMembersInts(key string) (reply []int, err error) {
	conn := c.pool.Get()
	defer conn.Close()

	return redis.Ints(conn.Do("SMEMBERS", key))
}

//SMembersStrings return all string members of key set
func (c *Client) SMembersStrings(key string) (reply []string, err error) {
	conn := c.pool.Get()
	defer conn.Close()

	return redis.Strings(conn.Do("SMEMBERS", key))
}

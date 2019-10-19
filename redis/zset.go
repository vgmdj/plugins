package redis

import (
	"github.com/gomodule/redigo/redis"
)

// ZADD add score and member to sorted set
// if already exist, return count = 0 and update member's score
// else return count = 1
func (c *Client) ZAdd(key string, score int, member string) (count int64, err error) {
	conn := c.pool.Get()
	defer conn.Close()

	return redis.Int64(conn.Do("ZADD", key, score, member))

}

// ZAddWithIncr incr score to member
// return count = origin score + incr score
// eg. already exist member = a, score = 10
//     zdd [key] incr 10 a
// return count = 20
// now member = a , score  = 20
func (c *Client) ZAddWithIncr(key string, score int, member string) (resultScore int64, err error) {
	conn := c.pool.Get()
	defer conn.Close()

	return redis.Int64(conn.Do("ZADD", key, "INCR", score, member))

}

// ZRange list members from start to end
func (c *Client) ZRange(key string, start, end int) ([]string, error) {
	conn := c.pool.Get()
	defer conn.Close()

	return redis.Strings(conn.Do("ZRANGE", key, start, end))

}

// ZRange list members from start to end
func (c *Client) ZRangeWithScore(key string, start, end int) ([]string, error) {
	conn := c.pool.Get()
	defer conn.Close()

	return redis.Strings(conn.Do("ZRANGE", key, start, end, "WITHSCORES"))

}

// ZRank return the index of the member
func (c *Client) ZRank(key string, member string) (int64, error) {
	conn := c.pool.Get()
	defer conn.Close()

	return redis.Int64(conn.Do("ZRANK", key, member))

}

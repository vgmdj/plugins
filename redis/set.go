package redis

import (
	"github.com/gomodule/redigo/redis"
	"github.com/vgmdj/utils/logger"
)

//SAdd add key-value and return the count of success
//notice: if you want use slice or array , you can only use int, int64, float64, string slice or array
//SAdd(key, 1,2,3,4,5,6)
//SAdd(key, []int{1,2,3,4,5,6})
//SAdd(key, []string{"1","2","3","4","5","6"})
func (c *Client) SAdd(key string, args ...interface{}) (count int64, err error) {
	conn := c.pool.Get()
	defer conn.Close()

	for _, v := range convertArgs(args[:]...) {
		success, err := redis.Int64(conn.Do("SADD", key, v))
		if err != nil {
			return count, err
		}
		count += success
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

//SMembers return all members of key set
func (c *Client) SIsMember(key string, member string) (bool, error) {
	conn := c.pool.Get()
	defer conn.Close()

	return redis.Bool(conn.Do("SISMEMBER", key, member))

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

func convertArgs(args ...interface{}) []interface{} {
	var params []interface{}

	if len(args) == 0 {
		return params
	}

	switch args[0].(type) {
	case []string:
		for _, v := range args[0].([]string) {
			params = append(params, v)
		}

	case []int:
		for _, v := range args[0].([]int) {
			params = append(params, v)
		}

	case []float64:
		for _, v := range args[0].([]float64) {
			params = append(params, v)
		}

	case []int64:
		for _, v := range args[0].([]int64) {
			params = append(params, v)
		}

	default:
		for _, v := range args {
			params = append(params, v)
		}
	}

	return params
}

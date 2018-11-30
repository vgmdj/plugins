package redis

import (
	"github.com/garyburd/redigo/redis"
	"github.com/vgmdj/utils/logger"
)

func Set(key, value interface{}) (err error) {
	c := redisPool.Get()
	defer c.Close()

	if _, err = c.Do("SET", key, value); err != nil {
		logger.Error(err.Error())
		return
	}

	return
}

func GetString(key string) (string, bool) {
	c := redisPool.Get()
	defer c.Close()
	count, _ := redis.Int(c.Do("EXISTS", key))
	if count == 0 {
		return "", false
	} else {
		n, _ := redis.String(c.Do("GET", key))
		return n, true
	}
}

func GetInt(key string) (int, bool) {
	c := redisPool.Get()
	defer c.Close()
	count, _ := redis.Int(c.Do("EXISTS", key))
	if count == 0 {
		return 0, false
	} else {
		n, _ := redis.Int(c.Do("GET", key))
		return n, true
	}
}

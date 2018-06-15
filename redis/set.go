package redis

import (
	"github.com/garyburd/redigo/redis"
	"log"
)

//SADD return count of the last one
func SAdd(key string, args ...interface{}) (count int64, err error) {
	c := redisPool.Get()
	defer c.Close()

	for _, v := range args {
		if count, err = redis.Int64(c.Do("SADD", key, v)); err != nil {
			log.Println(err.Error())
			return
		}
	}

	return
}

func SCard(key string) (count int64, err error) {
	c := redisPool.Get()
	defer c.Close()

	if count, err = redis.Int64(c.Do("SCARD", key)); err != nil {
		log.Println(err.Error())
		return
	}

	return
}

func SRem(key string, member interface{}) (err error) {
	c := redisPool.Get()
	defer c.Close()

	_, err = c.Do("SREM", key, member)
	if err != nil {
		log.Println(err.Error())
		return
	}

	return
}

func SMembersInt(key string) (reply []int, err error) {
	c := redisPool.Get()
	defer c.Close()

	if reply, err = redis.Ints(c.Do("SMEMBERS", key)); err != nil {
		log.Println(err.Error())
		return
	}
	return
}

func SMembersString(key string) (reply []string, err error) {
	c := redisPool.Get()
	defer c.Close()

	if reply, err = redis.Strings(c.Do("SMEMBERS", key)); err != nil {
		log.Println(err.Error())
		return
	}

	return
}

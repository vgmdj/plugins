package redis

import (
	"github.com/garyburd/redigo/redis"
	"github.com/vgmdj/utils/logger"
)

//TODO details

func HSet(key, itemKey string, item interface{}) (err error) {
	c := redisPool.Get()
	defer c.Close()

	if _, err = c.Do("HSET", key, itemKey, item); err != nil {
		logger.Error(err.Error())
		return
	}

	return
}

func HExists(key, itemKey string) bool {
	c := redisPool.Get()
	defer c.Close()
	count, _ := redis.Int(c.Do("HEXISTS", key, itemKey))
	if count == 0 {
		return false
	}
	return true
}

func HGet(userKey, itemKey string) (interface{}, bool) {
	c := redisPool.Get()
	defer c.Close()
	count, _ := redis.Int(c.Do("HEXISTS", userKey, itemKey))
	if count == 0 {
		return nil, false
	} else {
		res, _ := redis.Values(c.Do("HGET", userKey, itemKey))

		return res, true
	}
}

func HGetBool(userKey, itemKey string) (bool, bool) {
	c := redisPool.Get()
	defer c.Close()
	count, _ := redis.Int(c.Do("HEXISTS", userKey, itemKey))
	if count == 0 {
		return false, false
	} else {
		n, _ := redis.Bool(c.Do("HGET", userKey, itemKey))
		return n, true
	}
}

func HGetBytes(userKey, itemKey string) ([]byte, bool) {
	c := redisPool.Get()
	defer c.Close()
	count, _ := redis.Int(c.Do("HEXISTS", userKey, itemKey))
	if count == 0 {
		return nil, false
	} else {
		n, _ := redis.Bytes(c.Do("HGET", userKey, itemKey))
		return n, true
	}
}

func HGetFloat64(userKey, itemKey string) (float64, bool) {
	c := redisPool.Get()
	defer c.Close()
	count, _ := redis.Int(c.Do("HEXISTS", userKey, itemKey))
	if count == 0 {
		return 0, false
	} else {
		n, _ := redis.Float64(c.Do("HGET", userKey, itemKey))
		return n, true
	}
}

func HGetInt(userKey, itemKey string) (int, bool) {
	c := redisPool.Get()
	defer c.Close()
	count, _ := redis.Int(c.Do("HEXISTS", userKey, itemKey))
	if count == 0 {
		return 0, false
	} else {
		n, _ := redis.Int(c.Do("HGET", userKey, itemKey))
		return n, true
	}
}

func HGetInt64(userKey, itemKey string) (int64, bool) {
	c := redisPool.Get()
	defer c.Close()
	count, _ := redis.Int(c.Do("HEXISTS", userKey, itemKey))
	if count == 0 {
		return 0, false
	} else {
		n, _ := redis.Int64(c.Do("HGET", userKey, itemKey))
		return n, true
	}
}

func HGetIntMap(userKey, itemKey string) (map[string]int, bool) {
	c := redisPool.Get()
	defer c.Close()
	count, _ := redis.Int(c.Do("HEXISTS", userKey, itemKey))
	if count == 0 {
		return nil, false
	} else {
		n, _ := redis.IntMap(c.Do("HGET", userKey, itemKey))
		return n, true
	}
}

func HGetInt64Map(userKey, itemKey string) (map[string]int64, bool) {
	c := redisPool.Get()
	defer c.Close()
	count, _ := redis.Int(c.Do("HEXISTS", userKey, itemKey))
	if count == 0 {
		return nil, false
	} else {
		n, _ := redis.Int64Map(c.Do("HGET", userKey, itemKey))
		return n, true
	}
}

func HGetInts(userKey, itemKey string) ([]int, bool) {
	c := redisPool.Get()
	defer c.Close()
	count, _ := redis.Int(c.Do("HEXISTS", userKey, itemKey))
	if count == 0 {
		return nil, false
	} else {
		n, _ := redis.Ints(c.Do("HGET", userKey, itemKey))
		return n, true
	}
}

func HGetString(userKey, itemKey string) (string, bool) {
	c := redisPool.Get()
	defer c.Close()
	count, _ := redis.Int(c.Do("HEXISTS", userKey, itemKey))
	if count == 0 {
		return "", false
	} else {
		n, _ := redis.String(c.Do("HGET", userKey, itemKey))
		return n, true
	}
}

func HGetStrings(userKey, itemKey string) ([]string, bool) {
	c := redisPool.Get()
	defer c.Close()
	count, _ := redis.Int(c.Do("HEXISTS", userKey, itemKey))
	if count == 0 {
		return nil, false
	} else {
		n, _ := redis.Strings(c.Do("HGET", userKey, itemKey))
		return n, true
	}
}

func HGetStringMap(userKey, itemKey string) (map[string]string, bool) {
	c := redisPool.Get()
	defer c.Close()
	count, _ := redis.Int(c.Do("HEXISTS", userKey, itemKey))
	if count == 0 {
		return nil, false
	} else {
		n, _ := redis.StringMap(c.Do("HGET", userKey, itemKey))
		return n, true
	}
}

func HGetUint64(userKey, itemKey string) (uint64, bool) {
	c := redisPool.Get()
	defer c.Close()
	count, _ := redis.Int(c.Do("HEXISTS", userKey, itemKey))
	if count == 0 {
		return 0, false
	} else {
		n, _ := redis.Uint64(c.Do("HGET", userKey, itemKey))
		return n, true
	}
}

func HDel(userKey, itemKey string) (err error) {
	c := redisPool.Get()
	defer c.Close()
	_, err = c.Do("HDEL", userKey, itemKey)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	return
}

package redis

import (
	"github.com/garyburd/redigo/redis"
	"log"
	"sync"
)

//TODO

var (
	rds  *redisConf
	once sync.Once

	redisPool   *redis.Pool
	PoolMaxIdle = 10
)

type redisConf struct {
	address string
	pwd     string
	db      int64
}

func NewRedis(addr, pwd string, db int64) (err error) {
	once.Do(func() {
		rds = &redisConf{
			address: addr,
			pwd:     pwd,
			db:      db,
		}

		err := rds.init()
		if err != nil {
			return
		}
	})

	return
}

func (r *redisConf) init() error {
	redisPool = redis.NewPool(func() (redis.Conn, error) {
		c, err := redis.Dial("tcp", r.address)
		if err != nil {
			log.Println("--Redis--Connect redis fail:" + err.Error())
			return nil, err
		}
		if len(r.pwd) > 0 {
			if _, err := c.Do("AUTH", r.pwd); err != nil {
				c.Close()
				log.Println("--Redis--Auth redis fail:" + err.Error())
				return nil, err
			}
		}
		if _, err := c.Do("SELECT", r.db); err != nil {
			c.Close()
			log.Println("--Redis--Select redis db fail:" + err.Error())
			return nil, err
		}
		return c, nil
	}, PoolMaxIdle)
	return nil
}

func Store(key, itemKey string, item interface{}) {
	c := redisPool.Get()
	defer c.Close()

	if _, err := c.Do("HSET", key, itemKey, item); err != nil {
		log.Println(err.Error())
	}
}

func Expire(key string, seconds int) {
	c := redisPool.Get()
	defer c.Close()

	if _, err := c.Do("EXPIRE", key, seconds); err != nil {
		log.Println(err.Error())
	}
}

func Exists(key, itemKey string) bool {
	c := redisPool.Get()
	defer c.Close()
	count, _ := redis.Int(c.Do("HEXISTS", key, itemKey))
	if count == 0 {
		return false
	}
	return true
}

func Get(userKey, itemKey string) (interface{}, bool) {
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

func GetBool(userKey, itemKey string) (bool, bool) {
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

func GetBytes(userKey, itemKey string) ([]byte, bool) {
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

func GetFloat64(userKey, itemKey string) (float64, bool) {
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

func GetInt(userKey, itemKey string) (int, bool) {
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

func GetInt64(userKey, itemKey string) (int64, bool) {
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

func GetIntMap(userKey, itemKey string) (map[string]int, bool) {
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

func GetInt64Map(userKey, itemKey string) (map[string]int64, bool) {
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

func GetInts(userKey, itemKey string) ([]int, bool) {
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

func GetString(userKey, itemKey string) (string, bool) {
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

func GetStrings(userKey, itemKey string) ([]string, bool) {
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

func GetStringMap(userKey, itemKey string) (map[string]string, bool) {
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

func GetUint64(userKey, itemKey string) (uint64, bool) {
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

func RemoveItem(userKey, itemKey string) {
	c := redisPool.Get()
	defer c.Close()
	c.Do("HDEL", userKey, itemKey)
}

package redis

import (
	"github.com/garyburd/redigo/redis"
	"github.com/vgmdj/utils/logger"
	"sync"
)

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
	if addr == "" {
		addr = "127.0.0.1:6379"
	}

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

func IsOK() bool {
	if redisPool == nil {
		return false
	}

	return true
}

func (r *redisConf) init() error {
	redisPool = redis.NewPool(func() (redis.Conn, error) {
		c, err := redis.Dial("tcp", r.address)
		if err != nil {
			logger.Error("--Redis--Connect redis fail:" + err.Error())
			return nil, err
		}
		if len(r.pwd) > 0 {
			if _, err := c.Do("AUTH", r.pwd); err != nil {
				c.Close()
				logger.Error("--Redis--Auth redis fail:" + err.Error())
				return nil, err
			}
		}
		if _, err := c.Do("SELECT", r.db); err != nil {
			c.Close()
			logger.Error("--Redis--Select redis db fail:" + err.Error())
			return nil, err
		}
		return c, nil
	}, PoolMaxIdle)
	return nil
}

func Expire(key string, seconds int) (err error) {
	c := redisPool.Get()
	defer c.Close()

	if _, err = c.Do("EXPIRE", key, seconds); err != nil {
		logger.Error(err.Error())
		return
	}

	return
}

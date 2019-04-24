package redis

import (
	"testing"
	"time"
)

func TestClient_Expire(t *testing.T) {
	key := "TestExpire"

	c := testRedisClient()
	err := c.Set(key, 1)
	if err != nil {
		t.Error(err.Error())
		return
	}

	ok, _ := c.Exists(key)
	if !ok {
		t.Errorf("set key err")
		return
	}

	c.Expire(key, 3)

	time.Sleep(time.Second * 3)
	ok, _ = c.Exists(key)
	if ok {
		t.Errorf("expire key err")
		return
	}

	c.Set(key, 2)

	c.Expire(key, 0)
	ok, _ = c.Exists(key)
	if ok {
		t.Errorf("expire key err")
		return
	}

}

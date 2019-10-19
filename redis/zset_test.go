package redis

import (
	"testing"

	"github.com/vgmdj/utils/logger"
)

func TestZSet(t *testing.T) {
	c := testRedisClient()
	c.FlushAll()

	key := "TestSet"
	_, err := c.ZAdd(key, 16, "a")
	if err != nil {
		t.Error(err.Error())
		return
	}

	score, err := c.ZAddWithIncr(key, 4, "a")
	if err != nil {
		t.Error(err.Error())
		return
	}

	if score != 20 {
		t.Errorf("score should be 20, but get %d\n", score)
		return
	}

	for i := 0; i < 10; i++ {
		c.ZAdd(key, i+5, string('b'+i))
	}

	res, err := c.ZRangeWithScore(key, 0, 10)
	if err != nil {
		t.Error(err.Error())
		return
	}

	t.Log(res)

	index,err := c.ZRank(key,"f")
	if err != nil {
		logger.Error(err.Error())
		return
	}

	if index != 4{
		t.Errorf("index should be 4 , but get %d\n",index)
		return
	}

}

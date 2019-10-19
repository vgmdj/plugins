package redis

import (
	"strconv"
	"testing"
)

func TestSET(t *testing.T) {
	c := testRedisClient()

	key := "TestSet"
	members := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	err := c.DEL(key)
	if err != nil {
		t.Error(err.Error())
		return
	}

	count, err := c.SAdd(key, members)
	if err != nil {
		t.Error(err.Error())
		return
	}

	if count != 10 {
		t.Errorf("expected count :10, but got :%d", count)
		return
	}

	reply, err := c.SMembersInts(key)
	if err != nil {
		t.Error(err.Error())
		return

	}

	for k, v := range reply {
		if members[k] != v {
			t.Errorf("expected %d members is %d, but get %d", k, members[k], v)
			return
		}
	}

	count, err = c.SCard(key)
	if err != nil {
		t.Error(err.Error())
		return
	}

	if count != 10 {
		t.Errorf("expected count :10, but got :%d", count)
		return
	}

	exist, err := c.SIsMember(key, "9")
	if err != nil {
		t.Error(err.Error())
		return
	}

	if !exist {
		t.Errorf("expected true, but got false")
		return
	}

	c.SRem(key, "9")

	strReply, _ := c.SMembersStrings(key)
	for k, v := range strReply {
		if k == 8 {
			continue
		}

		if strconv.Itoa(members[k]) != v {
			t.Errorf("expected %d members is %d, but get %s", k, members[k], v)
			return
		}
	}

	exist, err = c.SIsMember(key, "9")
	if err != nil {
		t.Error(err.Error())
		return
	}

	if exist {
		t.Errorf("expected false, but got true")
		return
	}

}

func TestSetNX(t *testing.T) {
	c := testRedisClient()

	key := "TestSetNX"
	err := c.DEL(key)
	if err != nil {
		t.Error(err.Error())
		return

	}

	err = c.SetNX(key, "1", 30)
	if err != nil {
		t.Error(err.Error())
		return
	}

	err = c.SetNX(key, "1", 30)
	if err != ExistErr {
		t.Errorf("expected get exist err , but got : %s", err)
		return
	}

}

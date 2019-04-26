package redis

import "testing"

func TestSession(t *testing.T) {
	key := "TestSession"

	s := testRedisClient().NewSession()

	s.Begin()

	err := s.Exec("DEL", key)
	if err != nil {
		t.Error(err.Error())
		return
	}

	err = s.Exec("SET", key, 1)
	if err != nil {
		t.Error(err.Error())
		return
	}

	err = s.Exec("GET", key)
	if err != nil {
		t.Error(err.Error())
		return
	}

	err = s.Commit()
	if err != nil {
		t.Error(err.Error())
		return
	}

	for i := 0; i < 3; i++ {
		reply, err := s.Receive()
		if err != nil {
			t.Error(err.Error())
			return
		}
		t.Log(reply)
	}

}

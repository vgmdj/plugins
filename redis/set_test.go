package redis

import "testing"

func TestSET(t *testing.T) {
	NewRedis("", "", 0)

	SAdd("123", 1, 2, 3, 4, 5, 6, 7, 8, 9)

	members, _ := SMembersInt("123")
	t.Log(members)

	count, _ := SCard("123")
	t.Log(count)

	SRem("123", "9")

	strs, _ := SMembersString("123")
	t.Log(strs)

}

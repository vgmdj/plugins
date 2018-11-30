package redis

import "testing"

func TestSET(t *testing.T) {
	NewRedis("10.11.22.77:6379", "wangrui", 0)

	SAdd("123", 1, 2, 3, 4, 5, 6, 7, 8, 9)

	members, _ := SMembersInt("123")
	t.Log(members)

	count, _ := SCard("123")
	t.Log(count)

	SRem("123", "9")

	strs, _ := SMembersString("123")
	t.Log(strs)

}

func TestString(t *testing.T) {
	NewRedis("10.11.22.77:6379", "wangrui", 0)

	Set("123", 321)

	numbers, _ := GetInt("123")
	t.Log(numbers)

	str, _ := GetString("123")
	t.Log(str)

}


func TestSetNX(t *testing.T) {
	NewRedis("10.11.22.77:6379", "wangrui", 0)

	SetNX("124", 421,10)

	numbers, _ := GetInt("124")
	t.Log(numbers)


}
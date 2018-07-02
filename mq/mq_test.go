package mq

import "testing"

func TestNewMQ(t *testing.T) {
	conn := NewMQ(RabbitMQ, "", "", "", "")
	if conn == nil {
		t.Log("this is nil")
		return
	}

	t.Error(conn, "is not nil")

}

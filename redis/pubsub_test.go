package redis

import (
	"fmt"
	"testing"
	"time"

	"github.com/vgmdj/utils/logger"
)

func TestPubSub(t *testing.T) {
	channel := "TestChannel"
	c := testRedisClient()

	sub := func() {
		ps := c.NewPubSub()

		err := ps.Subscribe(channel)
		if err != nil {
			t.Error(err.Error())
			return
		}

		for {
			m, err := ps.Receive()
			if err != nil {
				t.Error(err.Error())
				return
			}

			logger.Info(m.Kind, m.Channel, string(m.Data))
		}

	}

	go sub()
	go sub()

	ps := c.NewPubSub()
	count := 0

	for count < 30 {
		values := fmt.Sprintf("prepare publish msg: send to %s nums %d", channel, count)

		nums, err := ps.Publish(channel, values)
		if err != nil {
			t.Error(err.Error())
			return
		}

		if nums != 2 {
			t.Errorf("expected nums is 2, but got %d\n", nums)
			return
		}

		count++

		time.Sleep(time.Second)
	}

}

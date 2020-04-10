package mns

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMultiConsume(t *testing.T) {
	assert := assert.New(t)

	c := &ConsumerConfig{
		URL:             "",
		AccessKeyID:     "",
		AccessKeySecret: "",

		Queue:       "",
		WaitSeconds: 1,
	}

	instance, err := NewConsumer(c)
	assert.Nil(err)

	consumer, ok := instance.(*consumer)
	assert.True(ok)

	summary := int64(0)

	proc := 1000
	wg := new(sync.WaitGroup)
	wg.Add(proc)

	consumer.LongPolling(10, false)
	consumer.LongPolling(10, true, 16)

	for i := 0; i < proc; i++ {
		go func() {
			defer wg.Done()

			for {
				resp, err := consumer.Next()
				if err == ErrMNSLibQuitWithCancel {
					break
				}

				if err != nil && err != ErrMNSPreClose {
					assert.Nil(err)
					break
				}

				//err = consumer.DeleteMessage(resp.ReceiptHandle)
				//assert.Nil(err)

				atomic.AddInt64(&summary, int64(len(resp)))
			}
		}()
	}

	go func() {
		time.Sleep(time.Second * 1)
		err = consumer.Close()
		if err != nil {
			assert.Nil(err)
		}
		t.Log("close done")
	}()

	wg.Wait()

	t.Log("consumed:", summary)
}

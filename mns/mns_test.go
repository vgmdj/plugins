package mns

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func getConf() *Config {
	return &Config{
		URL:             "",
		AccessKeyID:     "",
		AccessKeySecret: "",
	}
}

func TestSubsribePublish(t *testing.T) {
	assert := assert.New(t)

	{
		instance, err := NewMNSManager(getConf())
		assert.Nil(err)

		manager, ok := instance.(*mnsManager)
		assert.True(ok)

		err = manager.Subsribe("subsribeSyncXX", "queueSyncXX", "topicSyncXX", "")
		assert.Nil(err)
	}

	instance, err := NewMNS(getConf())
	assert.Nil(err)

	mns, ok := instance.(*mns)
	assert.True(ok)

	for k := 0; k < 10; k++ {
		ts := time.Now().Format(time.RFC3339)
		err = mns.Publish("topicSyncXX", "SyncXX",
			fmt.Sprintf("publish - %s - %d", ts, k))
		assert.Nil(err)
	}
}

func TestQueueSend(t *testing.T) {
	assert := assert.New(t)

	{
		instance, err := NewMNSManager(getConf())
		assert.Nil(err)

		manager, ok := instance.(*mnsManager)
		assert.True(ok)

		err = manager.CreateQueue("queueSyncYY", 0, 1024, 1000, 30, 0, 2)
		assert.Nil(err)
	}

	instance, err := NewMNS(getConf())
	assert.Nil(err)

	mns, ok := instance.(*mns)
	assert.True(ok)

	for k := 0; k < 10; k++ {
		ts := time.Now().Format(time.RFC3339)
		err = mns.Send("queueSyncYY", fmt.Sprintf("send - %s - %d", ts, k), DefaultMessagePriority)
		assert.Nil(err)
	}
}

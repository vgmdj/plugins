package mns

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMultiCreate(t *testing.T) {
	assert := assert.New(t)

	instance, err := NewMNSManager(getConf())
	assert.Nil(err)

	manager, ok := instance.(*mnsManager)
	assert.True(ok)

	topicExists := func(name string) {
		_, ok := manager.topics.Load(name)
		assert.True(ok)
	}

	queueExists := func(name string) {
		_, ok := manager.queues.Load(name)
		assert.True(ok)
	}

	err = manager.topicManager.DeleteTopic("topicA")
	assert.Nil(err)

	err = manager.queueManager.DeleteQueue("queueA")
	assert.Nil(err)

	err = manager.CreateTopic("topicA", 1024, true)
	assert.Nil(err)
	topicExists("topicA")

	err = manager.CreateTopic("topicA", 1024, true)
	assert.Nil(err)

	err = manager.CreateTopic("topicA", 2048, false)
	assert.Nil(err)

	err = manager.CreateQueue("queueA", 0, 1024, 1000, 30, 0, 2)
	assert.Nil(err)
	queueExists("queueA")

	err = manager.CreateQueue("queueA", 0, 1024, 1000, 30, 0, 2)
	assert.Nil(err)

	err = manager.CreateQueue("queueA", 0, 2048, 1000, 30, 0, 2)
	assert.Nil(err)
}

func TestFetch(t *testing.T) {
	assert := assert.New(t)

	instance, err := NewMNSManager(getConf())
	assert.Nil(err)

	manager, ok := instance.(*mnsManager)
	assert.True(ok)

	topicExists := func(name string) {
		_, ok := manager.topics.Load(name)
		assert.True(ok)
	}

	queueExists := func(name string) {
		_, ok := manager.queues.Load(name)
		assert.True(ok)
	}

	err = manager.CreateTopic("topicA", 1024, true)
	assert.Nil(err)
	topicExists("topicA")

	err = manager.CreateQueue("queueA", 0, 1024, 1000, 30, 0, 2)
	assert.Nil(err)
	queueExists("queueA")
}

func TestSubsribe(t *testing.T) {
	assert := assert.New(t)

	instance, err := NewMNSManager(getConf())
	assert.Nil(err)

	manager, ok := instance.(*mnsManager)
	assert.True(ok)

	topicExists := func(name string) {
		_, ok := manager.topics.Load(name)
		assert.True(ok)
	}

	queueExists := func(name string) {
		_, ok := manager.queues.Load(name)
		assert.True(ok)
	}

	err = manager.topicManager.DeleteTopic("topicA")
	assert.Nil(err)

	err = manager.queueManager.DeleteQueue("queueA")
	assert.Nil(err)

	err = manager.Subsribe("subsribeA", "queueA", "topicA", "")
	assert.Nil(err)
	topicExists("topicA")
	queueExists("queueA")

	err = manager.Subsribe("subsribeA", "queueA", "topicA", "")
	assert.Nil(err)
	topicExists("topicA")
	queueExists("queueA")
}

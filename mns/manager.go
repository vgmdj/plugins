package mns

import (
	"sync"

	alimns "github.com/aliyun/aliyun-mns-go-sdk"
	"github.com/pkg/errors"
)

const (
	// DefaultTopicMaxMessageSize 64k
	DefaultTopicMaxMessageSize = 65536
	// DefaultTopicLoggingEnabled disable
	DefaultTopicLoggingEnabled = false
	// DefaultQueueDelaySeconds no delay
	DefaultQueueDelaySeconds = 0
	// DefaultQueueMaxMessageSize 64k
	DefaultQueueMaxMessageSize = 65536
	// DefaultQueueMessageRetentionPeriod 4 days
	DefaultQueueMessageRetentionPeriod = 345600
	// DefaultQueueVisibilityTimeout message from inactive back to active
	DefaultQueueVisibilityTimeout = 30
	// DefaultQueuePollingWaitSeconds waitfor this time if no data in queue now
	DefaultQueuePollingWaitSeconds = 2
	// DefaultQueueSlices aliyun mns default
	DefaultQueueSlices = 2
)

// MNSManager supported operations of ali-mns-manager
type MNSManager interface {
	// CreateTopic try to fetch from aliyun if topic already exist
	// Deprecated: should operate in aliyun console
	CreateTopic(name string, maxMessageSize int32, loggingEnabled bool) error
	// CreateQueue try to fetch from aliyun if queue already exist
	// Deprecated: should operate in aliyun console
	CreateQueue(name string, delaySeconds, maxMessageSize, messageRetentionPeriod, visibilityTimeout, pollingWaitSeconds, slices int32) error
	// Subsribe try to create topic and queue if they are not exits, tag is used by filter
	Subsribe(subscriptionName, queueName, topicName, tagName string) error
}

type mnsManager struct {
	client       alimns.MNSClient
	topicManager alimns.AliTopicManager
	queueManager alimns.AliQueueManager
	topics       *sync.Map // string: alimns.AliMNSTopic
	queues       *sync.Map // string: alimns.AliMNSQueue
}

// NewMNSManager create a new instance of mns-manager
func NewMNSManager(conf *Config, accessToken ...string) (MNSManager, error) {
	if err := validateConf(conf); err != nil {
		return nil, err
	}

	token := ""
	if len(accessToken) != 0 && accessToken[0] != "" {
		token = accessToken[0]
	}

	manager := &mnsManager{
		topics: new(sync.Map),
		queues: new(sync.Map),
	}

	manager.client = alimns.NewAliMNSClientWithToken(conf.URL,
		conf.AccessKeyID, conf.AccessKeySecret, token)

	manager.topicManager = alimns.NewMNSTopicManager(manager.client)
	manager.queueManager = alimns.NewMNSQueueManager(manager.client)

	return manager, nil
}

func (m *mnsManager) CreateTopic(name string, maxMessageSize int32, loggingEnabled bool) error {
	if name == "" {
		return errors.New("name required")
	}

	if _, ok := m.topics.Load(name); ok {
		return nil
	}

	if _, err := m.topicManager.GetTopicAttributes(name); err == nil { // already exists
		m.topics.Store(name, alimns.NewMNSTopic(name, m.client))
		return nil
	}

	err := m.topicManager.CreateTopic(name, maxMessageSize, loggingEnabled)
	if err != nil && !alimns.ERR_MNS_TOPIC_ALREADY_EXIST.IsEqual(err) &&
		!alimns.ERR_MNS_TOPIC_ALREADY_EXIST_AND_HAVE_SAME_ATTR.IsEqual(err) {
		return errors.WithMessagef(err, "create topic err name: %s maxMessageSize: %d loggingEnabled: %v", name, maxMessageSize, loggingEnabled)
	}

	m.topics.Store(name, alimns.NewMNSTopic(name, m.client))
	return nil
}

func (m *mnsManager) CreateQueue(name string, delaySeconds, maxMessageSize, messageRetentionPeriod, visibilityTimeout, pollingWaitSeconds, slices int32) error {
	if name == "" {
		return errors.New("name required")
	}

	if _, ok := m.queues.Load(name); ok {
		return nil
	}

	if _, err := m.queueManager.GetQueueAttributes(name); err == nil { // already exists
		m.queues.Store(name, alimns.NewMNSQueue(name, m.client))
		return nil
	}

	err := m.queueManager.CreateQueue(name, delaySeconds, maxMessageSize, messageRetentionPeriod, visibilityTimeout, pollingWaitSeconds, slices)
	if err != nil && !alimns.ERR_MNS_QUEUE_ALREADY_EXIST.IsEqual(err) &&
		!alimns.ERR_MNS_QUEUE_ALREADY_EXIST_AND_HAVE_SAME_ATTR.IsEqual(err) {
		return errors.WithMessagef(err, "create queue err name: %s delaySeconds: %d maxMessageSize: %d messageRetentionPeriod: %d visibilityTimeout: %d pollingWaitSeconds: %d slices: %d",
			name, delaySeconds, maxMessageSize, messageRetentionPeriod, visibilityTimeout, pollingWaitSeconds, slices)
	}

	m.queues.Store(name, alimns.NewMNSQueue(name, m.client))
	return nil
}

func (m *mnsManager) Subsribe(subscriptionName, queueName, topicName, tagName string) error {
	if subscriptionName == "" {
		return errors.New("subscriptionName required")
	}
	if queueName == "" {
		return errors.New("queueName required")
	}
	if topicName == "" {
		return errors.New("topicName required")
	}
	// tag is optional

	if _, ok := m.queues.Load(queueName); !ok {
		if err := m.CreateQueue(queueName, DefaultQueueDelaySeconds, DefaultQueueMaxMessageSize,
			DefaultQueueMessageRetentionPeriod, DefaultQueueVisibilityTimeout, DefaultQueuePollingWaitSeconds, DefaultQueueSlices); err != nil {
			return err
		}
	}

	actual, ok := m.topics.Load(topicName)
	if !ok {
		if err := m.CreateTopic(topicName, DefaultTopicMaxMessageSize, DefaultTopicLoggingEnabled); err != nil {
			return err
		}
		actual, _ = m.topics.Load(topicName) // reload
	}
	topic := actual.(alimns.AliMNSTopic)

	err := topic.Subscribe(subscriptionName, alimns.MessageSubsribeRequest{
		Endpoint:            topic.GenerateQueueEndpoint(queueName),
		FilterTag:           tagName,
		NotifyStrategy:      alimns.BACKOFF_RETRY,
		NotifyContentFormat: alimns.SIMPLIFIED,
	})
	if err != nil && !alimns.ERR_MNS_SUBSCRIPTION_ALREADY_EXIST.IsEqual(err) &&
		!alimns.ERR_MNS_SUBSCRIPTION_ALREADY_EXIST_AND_HAVE_SAME_ATTR.IsEqual(err) {
		return errors.WithMessagef(err, "subscribe err desc: %s queue: %s topic: %s tag: %s", subscriptionName, queueName, topicName, tagName)
	}

	return nil
}

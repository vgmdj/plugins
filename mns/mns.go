package mns

import (
	"sync"

	alimns "github.com/aliyun/aliyun-mns-go-sdk"
	"github.com/pkg/errors"
)

var _ MNS = (*mns)(nil)

// MNS supported operations of ali-mns
type MNS interface {
	// Publish return err if topic not exist
	Publish(topic, tag, message string) error
	// Send return err if queue not exist
	Send(queue, message string, priority ...int64) error
}

// Config of aliyun mns
type Config struct {
	// URL mns endpoint
	URL string
	// AccessKeyID user's keyid
	AccessKeyID string
	// AccessKeySecret user's keysecret
	AccessKeySecret string

	accessToken     string
	channelCapacity int
}

type mns struct {
	client alimns.MNSClient
	topics map[string]alimns.AliMNSTopic
	queues map[string]alimns.AliMNSQueue
	mux    *sync.Mutex
}

func validateConf(conf *Config) error {
	if conf == nil {
		return errors.New("conf required")
	}

	if conf.URL == "" {
		return errors.New("conf.URL required")
	}
	if conf.AccessKeyID == "" {
		return errors.New("conf.AccessKeyID required")
	}
	if conf.AccessKeySecret == "" {
		return errors.New("conf.AccessKeySecret required")
	}

	return nil
}

// NewMNS create a new instance of mnslib
func NewMNS(conf *Config, accessToken ...string) (MNS, error) {
	if err := validateConf(conf); err != nil {
		return nil, err
	}

	token := ""
	if len(accessToken) != 0 && accessToken[0] != "" {
		token = accessToken[0]
	}

	return &mns{
		client: alimns.NewAliMNSClientWithToken(conf.URL,
			conf.AccessKeyID, conf.AccessKeySecret, token),
		topics: make(map[string]alimns.AliMNSTopic),
		queues: make(map[string]alimns.AliMNSQueue),
		mux:    new(sync.Mutex),
	}, nil
}

func (m *mns) Publish(topic, tag, message string) error {
	if topic == "" {
		return errors.New("topic required")
	}
	if tag == "" {
		return errors.New("tag required")
	}
	if message == "" {
		return errors.New("message required")
	}

	return m.publish(topic, tag, message)
}

func (m *mns) publish(topicName, tag, body string) error {
	topic := func() alimns.AliMNSTopic {
		m.mux.Lock()
		defer m.mux.Unlock()

		topic, ok := m.topics[topicName]
		if !ok {
			// just create a topic instance, if the topic not exist err will returned
			topic = alimns.NewMNSTopic(topicName, m.client)
			m.topics[topicName] = topic
		}

		return topic
	}()

	_, err := topic.PublishMessage(alimns.MessagePublishRequest{
		MessageTag:  tag,
		MessageBody: body,
	})
	if err != nil {
		return errors.WithMessage(err, "publish message err")
	}

	return nil
}

func (m *mns) Send(queue, message string, priority ...int64) error {
	if queue == "" {
		return errors.New("queue required")
	}
	if message == "" {
		return errors.New("message required")
	}

	var defaultPriority int64 = DefaultMessagePriority
	if len(priority) != 0 && priority[0] >= MinMessagePriority &&
		priority[0] <= MaxMessagePriority {
		defaultPriority = priority[0]
	}

	return m.send(queue, message, defaultPriority)
}

func (m *mns) send(queueName, body string, priority int64) error {
	queue := func() alimns.AliMNSQueue {
		m.mux.Lock()
		defer m.mux.Unlock()

		queue, ok := m.queues[queueName]
		if !ok {
			// just create a queue instance, if the queue not exist err will returned
			queue = alimns.NewMNSQueue(queueName, m.client)
			m.queues[queueName] = queue
		}

		return queue
	}()

	_, err := queue.SendMessage(alimns.MessageSendRequest{
		MessageBody: body,
		Priority:    priority,
	})
	if err != nil {
		return errors.WithMessage(err, "send message err")
	}

	return nil
}

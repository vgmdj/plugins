package mns

import (
	"context"
	"fmt"
	"io"
	"sync"
	"time"

	"go.uber.org/atomic"

	alimns "github.com/aliyun/aliyun-mns-go-sdk"
	"github.com/pkg/errors"
)

var _ Consumer = (*consumer)(nil)
var ErrMNSLibQuitWithCancel = io.EOF
var ErrMNSPreClose = fmt.Errorf("received close signal but still has msg in channel")

const (
	// DefaultConsumeWaitSeconds if aliyun return no data consumer will wait for this time and try next
	DefaultConsumeWaitSeconds = 1
	MaxConsumeWaitSeconds     = 30

	DefaultBatchNum int32 = 16
)

// Consumer consume message from queue
type Consumer interface {
	// Close consumer
	Close() error
	// LongPolling long polling fetch msg
	LongPolling(procs uint16, withBatch bool, batchNum ...uint8)
	// Next this will block if no data got, err maybe io.EOF if mnslib closed
	Next() (resp []alimns.MessageReceiveResponse, err error)
	// InvisibleMessage invisible message for some times
	InvisibleMessage(handle string, ttl time.Duration) error
	// DeleteMessage delete from queue
	DeleteMessage(handle string) error
}

type ConsumerConfig struct {
	// URL mns endpoint
	URL string
	// AccessKeyID user's keyid
	AccessKeyID string
	// AccessKeySecret user's keysecret
	AccessKeySecret string

	AccessToken string
	Queue       string
	WaitSeconds uint8
}

type consumer struct {
	ctx         context.Context
	cancel      context.CancelFunc
	queue       alimns.AliMNSQueue
	waitSeconds int64
	respCh      chan alimns.MessageReceiveResponse
	errCh       chan error
	batchCh     chan alimns.BatchMessageReceiveResponse
	wg          sync.WaitGroup
	task        chan bool
	taskNum     *atomic.Int32
	preClose    bool
	closeOnce   sync.Once
}

func validateConsumerConf(conf *ConsumerConfig) error {
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
	if conf.Queue == "" {
		return errors.New("queue required")
	}

	return nil
}

// NewConsumer create a new instance of consumer
func NewConsumer(conf *ConsumerConfig) (Consumer, error) {
	if err := validateConsumerConf(conf); err != nil {
		return nil, err
	}

	client := alimns.NewAliMNSClientWithToken(conf.URL,
		conf.AccessKeyID, conf.AccessKeySecret, conf.AccessToken)

	ctx, cancel := context.WithCancel(context.Background())
	return &consumer{
		ctx:         ctx,
		cancel:      cancel,
		queue:       alimns.NewMNSQueue(conf.Queue, client), // fetch queue from aliyun
		waitSeconds: setWaitSeconds(conf.WaitSeconds),
		respCh:      make(chan alimns.MessageReceiveResponse),
		errCh:       make(chan error),
		batchCh:     make(chan alimns.BatchMessageReceiveResponse),
		wg:          sync.WaitGroup{},
		task:        make(chan bool),
		taskNum:     atomic.NewInt32(0),
		preClose:    false,
		closeOnce:   sync.Once{},
	}, nil
}

func setWaitSeconds(n uint8) int64 {
	if n <= 0 {
		return DefaultConsumeWaitSeconds
	}

	if n > 30 {
		return MaxConsumeWaitSeconds
	}

	return int64(n)
}

// Close handle the rest channel msg then quit
// the Next() will return the ErrMnsQuitWithCancel
func (c *consumer) Close() error {
	var err = errors.New("has closed")
	c.preClose = true

	c.closeOnce.Do(func() {
		c.cancel()

		c.wg.Wait() // wait for all long polling exit
		close(c.batchCh)
		close(c.respCh)
		close(c.errCh)

		for i := int32(0); i < c.taskNum.Load(); i++ {
			<-c.task
		}

		close(c.task)

		err = nil
	})

	return err

}

// LongPolling fetch the msg to it's own buffer chan
func (c *consumer) LongPolling(procs uint16, withBatch bool, batchNum ...uint8) {
	if c.preClose {
		return
	}

	bn := DefaultBatchNum
	if len(batchNum) > 0 && int32(batchNum[0]) < DefaultBatchNum {
		bn = int32(batchNum[0])
	}

	for i := uint16(0); i < procs; i++ {
		go func() {
			// queue receive message or cancel context
			c.wg.Add(1)
			defer c.wg.Done()

			for {
				select {
				case <-c.ctx.Done():
					return

				case <-c.task:
					c.taskNum.Dec()
					if withBatch {
						c.queue.BatchReceiveMessage(c.batchCh, c.errCh, bn, c.waitSeconds)
						return
					}

					c.queue.ReceiveMessage(c.respCh, c.errCh, c.waitSeconds)

				}
			}

		}()
	}
}

func (c *consumer) result() (resp []alimns.MessageReceiveResponse, err error) {
	var (
		ok    bool
		batch alimns.BatchMessageReceiveResponse
		msg   alimns.MessageReceiveResponse
	)

	select {
	case msg, ok = <-c.respCh:
		resp = append(resp, msg)

	case batch, ok = <-c.batchCh:
		resp = batch.Messages

	case err, ok = <-c.errCh:

	}

	// channel closed
	if !ok {
		err = ErrMNSLibQuitWithCancel
	}

	return

}

func (c *consumer) Next() (resp []alimns.MessageReceiveResponse, err error) {
	for {
		if !c.preClose {
			c.taskNum.Inc()
			c.task <- true
		}

		resp, err = c.result()
		if alimns.ERR_MNS_MESSAGE_NOT_EXIST.IsEqual(err) {
			// queue is empty now and will retry next
			continue
		}

		if err != nil && err != ErrMNSLibQuitWithCancel {
			err = errors.WithMessage(err, "receive from queue err")
		}

		if c.preClose && err == nil {
			err = ErrMNSPreClose
		}

		return

	}
}

func (c *consumer) InvisibleMessage(handle string, ttl time.Duration) error {
	if handle != "" {
		_, err := c.queue.ChangeMessageVisibility(handle, int64(ttl))
		return err
	}

	return errors.New("handle required")
}

func (c *consumer) DeleteMessage(handle string) error {
	if handle != "" {
		return c.queue.DeleteMessage(handle)
	}

	return errors.New("handle required")
}

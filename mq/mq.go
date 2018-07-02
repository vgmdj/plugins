package mq

import (
	"github.com/vgmdj/plugins/mq/rabbitmq"
	"github.com/vgmdj/utils/logger"
)

type MQType int

const (
	RabbitMQ MQType = iota
)

type MQ interface {
}

func NewMQ(mt MQType, server, vhost, username, pwd string) MQ {
	var (
		mq  MQ
		err error
	)

	switch mt {
	default:
		return nil

	case RabbitMQ:
		mq, err = rabbitmq.NewRabbit(server, vhost, username, pwd)
	}

	if err != nil {
		logger.Error(err.Error())
		return nil
	}

	return mq

}

//default
func SetLogLevelDebug() {
	logger.SetLevel(logger.LevelDebug)
}

func SetLogLevelError() {
	logger.SetLevel(logger.LevelError)
}

func SetLogLevelEmergency() {
	logger.SetLevel(logger.LevelEmergency)
}

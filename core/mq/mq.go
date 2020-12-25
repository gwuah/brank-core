package mq

import "brank/core"

type MQ interface {
	Publish(topic string, msg []byte) error
	Subscribe(topics []string) chan []byte
	Close()
}

func NewMQ(c *core.Config) (MQ, error) {
	return NewKafkaMQ(c)
}
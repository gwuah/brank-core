package internal

import (
	"log"

	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

type EventStore interface {
	Publish(topic string, msg []byte) error
	Subscribe() chan []byte
	Close()
}

type messaging struct {
	topics []string
	broker string
	c      *kafka.Consumer
	p      *kafka.Producer
}

func NewEventStore(topics []string, broker, group string) EventStore {

	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"metadata.broker.list": broker,
	})

	if err != nil {
		panic(err)
	}

	go func() {
		for e := range producer.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					log.Printf("Delivery failed: %v\n", ev.TopicPartition)
				}
			}
		}
	}()

	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"metadata.broker.list": broker,
		"group.id":             group,
		"auto.offset.reset":    "earliest",
	})

	if err != nil {
		panic(err)
	}

	return &messaging{
		c:      consumer,
		p:      producer,
		topics: topics,
		broker: broker,
	}
}

func (m *messaging) Publish(topic string, msg []byte) error {

	err := m.p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          msg,
	}, nil)

	m.p.Flush(15 * 1000)

	return err
}

func (m *messaging) Subscribe() chan []byte {
	stream := make(chan []byte, 100)
	m.c.SubscribeTopics(m.topics, nil)

	go func() {
		for {
			msg, err := m.c.ReadMessage(-1)
			if err == nil {
				stream <- msg.Value
			} else {
				log.Printf("consumer error: %v (%v)\n", err, msg)
			}
		}
	}()

	return stream

}

func (m *messaging) Close() {
	defer m.p.Close()
}

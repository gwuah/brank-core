package internal

import (
	"log"

	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

type EventStore interface {
	Publish(topic string, msg []byte) error
	Subscribe(topics []string) chan []byte
	Close()
}

type messaging struct {
	c   *kafka.Consumer
	p   *kafka.Producer
	cfg *Config
}

func getProducerConfig(config *Config) *kafka.ConfigMap {
	if config.ENVIRONMENT == Development {
		return &kafka.ConfigMap{
			"metadata.broker.list": "localhost",
		}
	} else {
		return &kafka.ConfigMap{
			"metadata.broker.list":            config.CLOUDKARAFKA_BROKERS,
			"security.protocol":               "SASL_SSL",
			"sasl.mechanisms":                 "SCRAM-SHA-256",
			"sasl.username":                   config.CLOUDKARAFKA_USERNAME,
			"sasl.password":                   config.CLOUDKARAFKA_PASSWORD,
			"group.id":                        config.KAFKA_GROUP_ID,
			"go.events.channel.enable":        true,
			"go.application.rebalance.enable": true,
			"default.topic.config":            kafka.ConfigMap{"auto.offset.reset": "earliest"},
		}
	}
}

func getConsumerConfig(config *Config) *kafka.ConfigMap {
	if config.ENVIRONMENT == Development {
		return &kafka.ConfigMap{
			"metadata.broker.list": "localhost",
			"group.id":             config.KAFKA_GROUP_ID,
			"auto.offset.reset":    "earliest",
		}
	} else {
		return &kafka.ConfigMap{
			"metadata.broker.list":            config.CLOUDKARAFKA_BROKERS,
			"security.protocol":               "SASL_SSL",
			"sasl.mechanisms":                 "SCRAM-SHA-256",
			"sasl.username":                   config.CLOUDKARAFKA_USERNAME,
			"sasl.password":                   config.CLOUDKARAFKA_PASSWORD,
			"group.id":                        config.KAFKA_GROUP_ID,
			"go.events.channel.enable":        true,
			"go.application.rebalance.enable": true,
			"default.topic.config":            kafka.ConfigMap{"auto.offset.reset": "earliest"},
			"auto.offset.reset":               "earliest",
		}
	}
}

func NewEventStore(config *Config) EventStore {

	producer, err := kafka.NewProducer(getProducerConfig(config))

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

	consumer, err := kafka.NewConsumer(getConsumerConfig(config))

	if err != nil {
		panic(err)
	}

	return &messaging{
		c:   consumer,
		p:   producer,
		cfg: config,
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

func (m *messaging) Subscribe(topics []string) chan []byte {
	stream := make(chan []byte, 100)
	m.c.SubscribeTopics(topics, nil)

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

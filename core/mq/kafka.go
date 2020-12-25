package mq

import (
	"brank/core"
	"log"

	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

type kafkaMQ struct {
	c   *kafka.Consumer
	p   *kafka.Producer
	cfg *core.Config
}

func getProducerConfig(config *core.Config) *kafka.ConfigMap {
	log.Println("environment", config.ENVIRONMENT)
	if config.ENVIRONMENT == core.Development {
		return &kafka.ConfigMap{
			"metadata.broker.list": "localhost",
		}
	} else {
		return &kafka.ConfigMap{
			"metadata.broker.list": config.CLOUDKARAFKA_BROKERS,
			"security.protocol":    "SASL_SSL",
			"sasl.mechanisms":      "SCRAM-SHA-256",
			"sasl.username":        config.CLOUDKARAFKA_USERNAME,
			"sasl.password":        config.CLOUDKARAFKA_PASSWORD,
			"ssl.ca.location":      "./cloudkarafka.ca",
		}
	}
}

func getConsumerConfig(config *core.Config) *kafka.ConfigMap {
	log.Println("environment", config.ENVIRONMENT)
	if config.ENVIRONMENT == core.Development {
		return &kafka.ConfigMap{
			"metadata.broker.list": "localhost",
			"group.id":             config.KAFKA_GROUP_ID,
			"auto.offset.reset":    "earliest",
		}
	} else {
		return &kafka.ConfigMap{
			"metadata.broker.list": config.CLOUDKARAFKA_BROKERS,
			"security.protocol":    "SASL_SSL",
			"sasl.mechanisms":      "SCRAM-SHA-256",
			"sasl.username":        config.CLOUDKARAFKA_USERNAME,
			"sasl.password":        config.CLOUDKARAFKA_PASSWORD,
			"group.id":             config.KAFKA_GROUP_ID,
			"default.topic.config": kafka.ConfigMap{"auto.offset.reset": "earliest"},
			"auto.offset.reset":    "earliest",
			"ssl.ca.location":      "./cloudkarafka.ca",
		}
	}
}

func NewKafkaMQ(config *core.Config) (MQ, error) {

	producer, err := kafka.NewProducer(getProducerConfig(config))

	if err != nil {
		return nil, err
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
		return nil, err
	}

	return &kafkaMQ{
		c:   consumer,
		p:   producer,
		cfg: config,
	}, nil
}

func (m *kafkaMQ) Publish(topic string, msg []byte) error {

	err := m.p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          msg,
	}, nil)

	m.p.Flush(15 * 1000)

	return err
}

func (m *kafkaMQ) Subscribe(topics []string) chan []byte {
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

func (m *kafkaMQ) Close() {
	defer m.p.Close()
}

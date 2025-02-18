package kafka

import (
	"errors"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"strings"
)

const (
	flushTimeout = 5000 // ms
)

var errUnknownType = errors.New("unknown even type")

type Producer struct {
	producer *kafka.Producer
}

func NewProducer(address []string) (*Producer, error) {
	config := &kafka.ConfigMap{
		"bootstrap.servers": strings.Join(address, ","),
	}
	p, err := kafka.NewProducer(config)
	if err != nil {
		return nil, fmt.Errorf("error with new producer: %w", err)
	}
	return &Producer{producer: p}, nil
}

func (p *Producer) Produce(message, topic string) error {
	kafkaMsg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Value: []byte(message),
		Key:   nil,
	}
	kafkaChan := make(chan kafka.Event)
	if err := p.producer.Produce(kafkaMsg, kafkaChan); err != nil {
		return err
	}
	e := <-kafkaChan
	switch ev := e.(type) {
	case *kafka.Message:
		return nil
	case kafka.Error:
		return ev
	default:
		return errUnknownType
	}
}

func (p *Producer) Close() {
	p.producer.Flush(flushTimeout)
	p.producer.Close()
}

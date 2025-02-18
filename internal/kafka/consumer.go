package kafka

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/labstack/gommon/log"
	"os"
	"strings"
	"time"
)

type Consumer struct {
	consumer *kafka.Consumer
	stop     bool
}

func NewConsumer(address []string, topic, consumerGroup string) (*Consumer, error) {
	cfg := &kafka.ConfigMap{
		"bootstrap.servers":        strings.Join(address, ","),
		"group.id":                 consumerGroup,
		"enable.auto.offset.store": false,
		"enable.auto.commit":       true,
	}
	c, err := kafka.NewConsumer(cfg)
	if err != nil {
		return nil, fmt.Errorf("Could not create new consumer: %v", err)
	}

	if err = c.Subscribe(topic, nil); err != nil {
		return nil, fmt.Errorf("Could not subscribe on topic: %v", err)
	}
	return &Consumer{consumer: c}, nil
}

func (c *Consumer) Start() {
	file, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Could not open log file: %v", err)
	}
	defer c.consumer.Close()
	defer file.Close()

	for {
		if c.stop {
			break
		}
		kafkaMsg, err := c.consumer.ReadMessage(-1)
		if err != nil {
			log.Error(err)
			continue
		}
		msg := fmt.Sprintf("Received message at %s: %s\n", time.Now().Format(time.RFC3339), string(kafkaMsg.Value))
		file.WriteString(msg)
	}
}

func (c *Consumer) Stop() error {
	c.stop = true
	return c.consumer.Close()
}

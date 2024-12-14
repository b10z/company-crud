package producer

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"time"
)

type Config struct {
	Server string
	Topic  string
}

type Produce interface {
	ProduceEvent(message []byte) error
}

type KafkaProducer struct {
	*kafka.Producer
	cfg Config
}

func New(conf Config) (*KafkaProducer, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": conf.Server})
	if err != nil {
		return nil, err
	}

	return &KafkaProducer{
		Producer: p,
		cfg:      conf,
	}, nil
}

func (kp *KafkaProducer) ProduceEvent(message []byte) error {
	err := kp.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &kp.cfg.Topic, Partition: kafka.PartitionAny},
		Value:          []byte(message),
	}, nil)

	if err != nil {
		if err.(kafka.Error).Code() == kafka.ErrQueueFull {
			time.Sleep(time.Second)
		}
		return fmt.Errorf("error while producing event: %w", err)
	}

	return nil
}

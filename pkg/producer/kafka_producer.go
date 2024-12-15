package producer

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"time"
)

type Config struct {
	Server string
	Acks   string
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
		"bootstrap.servers": conf.Server,
		"acks":              conf.Acks,
	})
	if err != nil {
		return nil, err
	}

	return &KafkaProducer{
		Producer: p,
		cfg:      conf,
	}, nil
}

func (kp *KafkaProducer) Stop() {
	kp.Flush(int(5 * time.Second))
	kp.Close()
}

func (kp *KafkaProducer) ProduceEvent(message []byte) error {
	go func() {
		for e := range kp.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()

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

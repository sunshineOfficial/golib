package gokafka

import (
	"github.com/segmentio/kafka-go"
	"github.com/sunshineOfficial/golib/goctx"
	"github.com/sunshineOfficial/golib/golog"
)

type Message kafka.Message

type Kafka interface {
	Producer(topicName string, options ...ProducerOption) Producer
	Consumer(log golog.Logger, getCtx goctx.ProvideWithCancel, options ...ConsumerOption) (Consumer, error)
}

type KafkaImpl struct {
	brokers []string
}

func NewKafka(brokers []string) *KafkaImpl {
	return &KafkaImpl{
		brokers: brokers,
	}
}

func (k *KafkaImpl) Producer(topicName string, options ...ProducerOption) Producer {
	return NewProducer(k.brokers, topicName, options...)
}

func (k *KafkaImpl) Consumer(log golog.Logger, getCtx goctx.ProvideWithCancel, options ...ConsumerOption) (Consumer, error) {
	return NewConsumer(log, getCtx, k.brokers, options...)
}

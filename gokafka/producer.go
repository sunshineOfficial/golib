package gokafka

import (
	"context"

	"github.com/segmentio/kafka-go"
	"github.com/sunshineOfficial/golib/gosync"
)

type Producer interface {
	Produce(ctx context.Context, message Message) error
	Close(ctx context.Context) error
}

type ProducerImpl struct {
	writer *kafka.Writer
}

func NewProducer(brokers []string, topic string, options ...ProducerOption) *ProducerImpl {
	return &ProducerImpl{
		writer: newKafkaWriter(brokers, topic, options),
	}
}

func (p *ProducerImpl) Produce(ctx context.Context, message Message) error {
	return p.writer.WriteMessages(ctx, kafka.Message(message))
}

func newKafkaWriter(brokers []string, topic string, options []ProducerOption) *kafka.Writer {
	opt := ProducerOptions{
		BatchSize:  -1,
		BatchBytes: -1,
	}

	for _, o := range options {
		opt = o(opt)
	}

	writer := &kafka.Writer{
		Addr:  kafka.TCP(brokers...),
		Topic: topic,
	}

	if opt.BatchSize > 0 {
		writer.BatchSize = opt.BatchSize
	}
	if opt.BatchBytes > 0 {
		writer.BatchBytes = opt.BatchBytes
	}
	writer.Async = opt.Async
	return writer
}

func (p *ProducerImpl) Close(ctx context.Context) error {
	return gosync.WaitContext(ctx, p.writer.Close)
}

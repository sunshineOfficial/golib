package gokafka

import (
	"context"

	"github.com/segmentio/kafka-go"
	"github.com/sunshineOfficial/golib/goctx"
	"github.com/sunshineOfficial/golib/golog"
	"github.com/sunshineOfficial/golib/gosync"
)

type Subscriber func(message Message, err error)

type Consumer interface {
	Consume(ctx context.Context) (Message, error)
	Subscribe(fn Subscriber)
	Close(ctx context.Context) error
}

type ConsumerImpl struct {
	reader   *kafka.Reader
	listener *listener
}

func NewConsumer(log golog.Logger, getContext goctx.ProvideWithCancel, brokers []string, options ...ConsumerOption) (*ConsumerImpl, error) {
	reader, err := newKafkaReader(brokers, options)
	if err != nil {
		return nil, err
	}

	consumer := &ConsumerImpl{
		reader: reader,
	}

	consumer.listener = newListener(log, func() (Message, error) {
		ctx, cancel := getContext()
		defer cancel()

		return consumer.Consume(ctx)
	})

	return consumer, nil
}

func newKafkaReader(brokers []string, options []ConsumerOption) (*kafka.Reader, error) {
	opt := ConsumerOptions{
		Partition:     -1,
		QueueCapacity: -1,
		MinBytes:      -1,
		MaxBytes:      -1,
		StartOffset:   0,
	}

	for _, o := range options {
		opt = o(opt)
	}

	cfg := kafka.ReaderConfig{
		Brokers: brokers,
	}

	if len(opt.GroupId) > 0 {
		cfg.GroupID = opt.GroupId
	}
	if len(opt.TopicName) > 0 {
		cfg.Topic = opt.TopicName
	}
	if opt.Partition > -1 {
		cfg.Partition = opt.Partition
	}
	if opt.QueueCapacity > -1 {
		cfg.QueueCapacity = opt.QueueCapacity
	}
	if opt.MinBytes > -1 {
		cfg.MinBytes = opt.MinBytes
	}
	if opt.MaxBytes > -1 {
		cfg.MaxBytes = opt.MaxBytes
	}

	cfg.StartOffset = opt.StartOffset

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return kafka.NewReader(cfg), nil
}

func (c *ConsumerImpl) Consume(ctx context.Context) (Message, error) {
	message, err := c.reader.ReadMessage(ctx)

	return Message(message), err
}

func (c *ConsumerImpl) Subscribe(s Subscriber) {
	c.listener.add(s)
	c.listener.start()
}

func (c *ConsumerImpl) Close(ctx context.Context) error {
	return gosync.WaitContext(ctx, c.reader.Close)
}

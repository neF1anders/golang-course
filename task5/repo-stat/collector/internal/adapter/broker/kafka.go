package broker

import (
	"context"
	"log/slog"
	"repo-stat/collector/internal/usecase"
	"sync"
	"time"

	kafka "github.com/segmentio/kafka-go"
)

type Consumer struct {
	log     *slog.Logger
	reader  *kafka.Reader
	handler usecase.MessageHandler
	wg      sync.WaitGroup
	cancel  context.CancelFunc
}

func NewConsumer(brokers []string, groupID, topic string, handler usecase.MessageHandler, log *slog.Logger) (*Consumer, error) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        brokers,
		GroupID:        groupID,
		Topic:          topic,
		MinBytes:       1e3,
		MaxBytes:       1e6,
		MaxWait:        1 * time.Second,
		CommitInterval: 0,
	})
	return &Consumer{
		log:     log,
		reader:  r,
		handler: handler,
	}, nil
}
func (c *Consumer) Start(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	c.cancel = cancel

	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		for {
			if ctx.Err() != nil {
				return
			}
			msg, err := c.reader.FetchMessage(ctx)
			if err != nil {
				if ctx.Err() != nil {
					return
				}
				c.log.Error("failed to fetch message in kafka consumer", "error", err)
				continue
			}
			if err := c.handler.Handle(ctx, msg.Value); err != nil {
				c.log.Error("failed to handle the message", "error", err)
				continue
			}
			if err := c.reader.CommitMessages(ctx, msg); err != nil {
				c.log.Error("failed to commit message", "error", err)
				continue
			}
		}
	}()
	return nil
}
func (c *Consumer) Stop() error {
	if c.cancel != nil {
		c.cancel()
	}
	c.wg.Wait()
	return c.reader.Close()
}

type Producer struct {
	log    *slog.Logger
	writer *kafka.Writer
}

func NewProducer(brokers []string, topic string, log *slog.Logger) *Producer {
	w := &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Topic:        topic,
		Balancer:     &kafka.RoundRobin{},
		RequiredAcks: kafka.RequireAll,
		MaxAttempts:  5,
	}
	return &Producer{writer: w, log: log}
}
func (p *Producer) Publish(ctx context.Context, topic string, key, value []byte) error {
	msg := kafka.Message{
		Topic: topic,
		Key:   key,
		Value: value,
	}
	err := p.writer.WriteMessages(ctx, msg)
	if err != nil {
		p.log.Error("kafka producer error", "error", err)
	}
	return err
}
func (p *Producer) Close() error {
	return p.writer.Close()
}

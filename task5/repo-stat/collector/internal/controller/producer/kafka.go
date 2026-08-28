package producer

import (
	"context"
	"encoding/json"
	"repo-stat/collector/internal/adapter/broker"
	"repo-stat/collector/internal/domain"
)

type ResultPublisher struct {
	producer *broker.Producer
	topic    string
}

func NewResultPublisher(producer *broker.Producer, topic string) *ResultPublisher {
	return &ResultPublisher{
		producer: producer,
		topic:    topic,
	}
}
func (p *ResultPublisher) Publish(ctx context.Context, data []*domain.Repo) error {
	raw, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return p.producer.Publish(ctx, p.topic, nil, raw)
}

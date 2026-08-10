package usecase

import (
	"context"
	"repo-stat/api/internal/domain"
)

type Ping struct {
	processor_pinger  Pinger
	subscriber_pinger Pinger
}

func NewPing(processor_pinger Pinger, subscriber_pinger Pinger) *Ping {
	return &Ping{
		processor_pinger:  processor_pinger,
		subscriber_pinger: subscriber_pinger,
	}
}

func (u *Ping) Execute(ctx context.Context) (domain.PingStatus, domain.PingStatus) {
	processor_ping, _ := u.processor_pinger.Ping(ctx)
	subscriber_ping, _ := u.subscriber_pinger.Ping(ctx)
	return processor_ping, subscriber_ping
}

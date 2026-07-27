package consumer

import (
	"context"
	"log/slog"
	"repo-stat/collector/internal/usecase"
	"time"
)

type Scheduler struct {
	log       *slog.Logger
	collectUC *usecase.GetAndPublishUseCase
	interval  time.Duration
}

func NewScheduler(collectUC *usecase.GetAndPublishUseCase, interval time.Duration) *Scheduler {
	return &Scheduler{collectUC: collectUC, interval: interval}
}

func (s *Scheduler) Start(ctx context.Context) {
	ticker := time.NewTicker(s.interval)
	go func() {
		for {
			select {
			case <-ticker.C:
				if err := s.collectUC.Execute(ctx); err != nil {
					s.log.Error("scheduled collect failed", "error", err)
				}
			case <-ctx.Done():
				ticker.Stop()
				return
			}
		}
	}()
}

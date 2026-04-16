package usecase

import (
	"context"
	"repo-stat/subscribe/internal/domain"
)

type PingUseCase struct {
	pinger domain.Pinger
}

func NewPingerUseCase(pinger domain.Pinger) *PingUseCase {
	return &PingUseCase{pinger: pinger}
}
func (uc *PingUseCase) PingRepo(ctx context.Context, slug domain.Slug) error {
	return uc.pinger.PingRepo(ctx, slug)
}

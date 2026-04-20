package usecase

import (
	"context"
	domain "repo-stat/collector/internal/domain"
)

type GetSubInfoUseCase struct {
	fetcher   domain.Fetcher
	subscribe domain.Subscribe
}

func NewGetSubInfoUseCase(fetcher domain.Fetcher, subscribe domain.Subscribe) *GetSubInfoUseCase {
	return &GetSubInfoUseCase{fetcher: fetcher, subscribe: subscribe}
}
func (uc *GetSubInfoUseCase) Execute(ctx context.Context) ([]*domain.Repo, error) {
	slugs, err := uc.subscribe.Subinfo(ctx)
	if err != nil {
		return nil, err
	}
	res := make([]*domain.Repo, len(slugs))
	for i := 0; i < len(slugs); i++ {
		repo, err := uc.fetcher.Fetch(slugs[i].Owner, slugs[i].Repo)
		if err != nil {
			return res, err
		}
		res = append(res, repo)
	}
	return res, nil
}

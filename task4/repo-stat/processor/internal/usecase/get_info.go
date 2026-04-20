package usecase

import (
	"context"
	"repo-stat/processor/internal/domain"
)

type Fetch struct {
	fetcher Fetcher
}

func NewFetch(fetcher Fetcher) *Fetch {
	return &Fetch{
		fetcher: fetcher,
	}
}

func (f *Fetch) GetInfo(ctx context.Context, owner, repo string) (domain.Repo, error) {
	return f.fetcher.GetInfo(ctx, owner, repo)
}

func (f *Fetch) GetSubInfo(ctx context.Context) ([]*domain.Repo, error) {
	return f.fetcher.GetSubInfo(ctx)
}

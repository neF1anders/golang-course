package usecase

import (
	"context"
	"repo-stat/api/internal/domain"
)

type Fetch struct {
	fetcher Fetcher
}

func NewFetch(fetcher Fetcher) *Fetch {
	return &Fetch{
		fetcher: fetcher,
	}
}

func (f *Fetch) GetSubInfo(ctx context.Context) ([]domain.Repo, error) {
	return f.fetcher.GetSubInfo(ctx)
}

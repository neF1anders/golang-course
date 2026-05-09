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

func (f *Fetch) Execute(ctx context.Context, owner, repo string) (domain.Repo, error) {
	return f.fetcher.GetInfo(ctx, owner, repo)
}

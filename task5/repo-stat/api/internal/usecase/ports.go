package usecase

import (
	"context"
	"repo-stat/api/internal/domain"
)

type Pinger interface {
	Ping(ctx context.Context) (domain.PingStatus, error)
}
type Fetcher interface {
	GetInfo(ctx context.Context, owner, repo string) (domain.Repo, error)
	GetSubInfo(ctx context.Context) ([]domain.Repo, error)
}
type Retriever interface {
	Subscribe(ctx context.Context, slug domain.Slug) error
	Unsubscribe(ctx context.Context, slug domain.Slug) error
	SubInfo(ctx context.Context) ([]domain.Slug, error)
}

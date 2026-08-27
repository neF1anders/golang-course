package usecase

import (
	"context"
	"repo-stat/processor/internal/domain"
)

/*
	type Fetcher interface {
		GetSubInfo(ctx context.Context) ([]*domain.Repo, error)
	}
*/
type DB interface {
	List(ctx context.Context) ([]domain.Repo, error)
	Push(ctx context.Context, repo domain.Repo) error
	Delete(ctx context.Context, repo domain.Repo) error
}
type Producer interface {
	Publish(ctx context.Context) error
	Close() error
}

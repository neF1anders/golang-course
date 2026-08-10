package usecase

import (
	"context"
	"repo-stat/processor/internal/domain"
)

type Fetcher interface {
	GetInfo(ctx context.Context, owner, repo string) (domain.Repo, error)
}

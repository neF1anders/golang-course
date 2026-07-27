package usecase

import (
	"context"
	domain "repo-stat/collector/internal/domain"
)

type Publisher interface {
	Publish(ctx context.Context, data []*domain.Repo) error
}
type MessageHandler interface {
	Handle(ctx context.Context, data []byte) error
}

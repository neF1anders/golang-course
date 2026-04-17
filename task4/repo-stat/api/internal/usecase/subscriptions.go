package usecase

import (
	"context"
	"repo-stat/api/internal/domain"
)

type RetrieverUseCase struct {
	db Retriever
}

func NewRetrieverUseCase(db Retriever) *RetrieverUseCase {
	return &RetrieverUseCase{
		db: db,
	}
}

func (u *RetrieverUseCase) Subscribe(ctx context.Context, slug domain.Slug) error {
	return u.db.Subscribe(ctx, slug)
}
func (u *RetrieverUseCase) Unsubscribe(ctx context.Context, slug domain.Slug) error {
	return u.db.Unsubscribe(ctx, slug)
}
func (u *RetrieverUseCase) SubInfo(ctx context.Context) ([]domain.Slug, error) {
	slugs, err := u.db.SubInfo(ctx)
	if err != nil {
		return nil, err
	}
	return slugs, nil
}

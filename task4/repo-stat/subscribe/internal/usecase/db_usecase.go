package usecase

import (
	"context"
	"fmt"
	"repo-stat/subscribe/internal/domain"
)

type DBUseCase struct {
	retriever domain.Retriever
}

func NewDBUseCase(retriever domain.Retriever) *DBUseCase {
	return &DBUseCase{retriever: retriever}
}
func (uc *DBUseCase) List(ctx context.Context) ([]domain.Slug, error) {
	return uc.retriever.List(ctx)
}
func (uc *DBUseCase) Push(ctx context.Context, slug domain.Slug) error {
	slugs, err := uc.retriever.List(ctx)
	if err != nil {
		return fmt.Errorf("failed to list subs: %v", err)
	}
	for _, el := range slugs {
		if el == slug {
			return fmt.Errorf("failed to push the sub %s/%s, dublication error", slug.Owner, slug.Repo)
		}
	}
	return uc.retriever.Push(ctx, slug)
}
func (uc *DBUseCase) Delete(ctx context.Context, slug domain.Slug) error {
	return uc.retriever.Delete(ctx, slug)
}

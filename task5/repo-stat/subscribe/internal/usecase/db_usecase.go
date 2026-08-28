package usecase

import (
	"context"
	"fmt"
	"repo-stat/subscribe/internal/domain"
)

type DBUseCase struct {
	retriever domain.Retriever
	pinger    *PingUseCase
}

func NewDBUseCase(retriever domain.Retriever, pinger *PingUseCase) *DBUseCase {
	return &DBUseCase{retriever: retriever, pinger: pinger}
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
		if el.Owner == slug.Owner && el.Repo == slug.Repo {
			return fmt.Errorf("failed to push the sub %s/%s: 409: dublication error", slug.Owner, slug.Repo)
		}
	}
	err = uc.pinger.PingRepo(ctx, slug)
	if err != nil {
		return fmt.Errorf("failed to push the sub %s/%s: %v", slug.Owner, slug.Repo, err)
	}
	return uc.retriever.Push(ctx, slug)
}
func (uc *DBUseCase) Delete(ctx context.Context, slug domain.Slug) error {
	slugs, err := uc.retriever.List(ctx)
	if err != nil {
		return fmt.Errorf("failed to list subs: %v", err)
	}
	var flag = false
	for _, el := range slugs {
		if el.Owner == slug.Owner && el.Repo == slug.Repo {
			flag = true
		}
	}
	if !flag {
		return fmt.Errorf("failed to delete the sub %s/%s: %v", slug.Owner, slug.Repo, "404: subscription not found")
	}
	return uc.retriever.Delete(ctx, slug)
}

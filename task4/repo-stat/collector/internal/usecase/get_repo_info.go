package usecase

import (
	domain "repo-stat/collector/internal/domain"
)

type GetRepoInfoUseCase struct {
	fetcher domain.Fetcher
}

func NewGetRepoInfoUseCase(fetcher domain.Fetcher) *GetRepoInfoUseCase {
	return &GetRepoInfoUseCase{fetcher: fetcher}
}
func (uc *GetRepoInfoUseCase) Execute(owner, repo string) (*domain.Repo, error) {
	return uc.fetcher.Fetch(owner, repo)
}

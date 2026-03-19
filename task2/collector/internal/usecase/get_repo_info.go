package usecase

import (
	"collector/internal/entity"
)

type GitHubRepoFetcher interface {
	Fetch(owner, repo string) (*entity.Repo, error)
}
type GetRepoInfoUseCase struct {
	fetcher GitHubRepoFetcher
}

func NewGetRepoInfoUseCase(fetcher GitHubRepoFetcher) *GetRepoInfoUseCase {
	return &GetRepoInfoUseCase{fetcher: fetcher}
}
func (uc *GetRepoInfoUseCase) Execute(owner, repo string) (*entity.Repo, error) {
	return uc.fetcher.Fetch(owner, repo)
}

package usecase

import (
	"context"
	"fmt"
	"repo-stat/processor/internal/domain"
)

type DBUseCase struct {
	retriever DB
}

func NewDBUseCase(db DB) *DBUseCase {
	return &DBUseCase{
		retriever: db,
	}
}

func (uc *DBUseCase) List(ctx context.Context) ([]domain.Repo, error) {
	return uc.retriever.List(ctx)
}
func (uc *DBUseCase) Update(ctx context.Context, repos []domain.Repo) error {
	repos_old, err := uc.retriever.List(ctx)
	if err != nil {
		return fmt.Errorf("failed to list repos from db: %v", err)
	}
	for _, el := range repos_old {
		if err := uc.retriever.Delete(ctx, el); err != nil {
			return fmt.Errorf("failed to delete repos into db: %v", err)
		}
	}
	for _, el := range repos {
		if err := uc.retriever.Push(ctx, el); err != nil {
			return fmt.Errorf("failed to push repos into db: %v", err)
		}
	}
	return nil
}

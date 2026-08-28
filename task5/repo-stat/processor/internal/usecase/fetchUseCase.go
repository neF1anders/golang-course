package usecase

import (
	"context"
	"fmt"
	"repo-stat/processor/internal/domain"
	"time"
)

type FetchUseCase struct {
	dbUC     *DBUseCase
	producer Producer
}

func NewFetchUC(dbUC *DBUseCase, producer Producer) *FetchUseCase {
	return &FetchUseCase{
		dbUC:     dbUC,
		producer: producer,
	}
}

func (f *FetchUseCase) GetSubInfo(ctx context.Context) ([]domain.Repo, error) {
	data, err := f.dbUC.List(ctx)
	if err != nil {
		return nil, err
	}
	if len(data) > 0 {
		return data, nil
	}
	if err = f.producer.Publish(ctx); err != nil {
		return nil, err
	}
	for i := 0; i < 5; i++ {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(500 * time.Millisecond):
			data, err := f.dbUC.List(ctx)
			if err == nil {
				return data, nil
			}
		}
	}
	return nil, fmt.Errorf("timeout waiting for subscription update")
}

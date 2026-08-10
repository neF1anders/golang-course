package domain

import (
	"context"
)

type Pinger interface {
	PingRepo(ctx context.Context, slug Slug) error
}
type Retriever interface {
	List(ctx context.Context) ([]Slug, error)
	Push(ctx context.Context, slug Slug) error
	Delete(ctx context.Context, slug Slug) error
}

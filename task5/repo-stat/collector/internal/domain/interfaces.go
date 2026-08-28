package domain

import "context"

type Fetcher interface {
	Fetch(owner, repo string) (*Repo, error)
}

type Subscribe interface {
	Subinfo(ctx context.Context) ([]*Slug, error)
}

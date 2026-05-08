package postgres

import (
	"context"
	"fmt"

	db "repo-stat/subscribe/db/sqlc"
	"repo-stat/subscribe/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresClient struct {
	pool *pgxpool.Pool
	q    *db.Queries
}

func NewPostgresClient(ctx context.Context, dsn string) (*PostgresClient, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to init connection pool, %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	return &PostgresClient{
		pool: pool,
		q:    db.New(pool),
	}, nil

}

func (r *PostgresClient) Close() {
	r.pool.Close()
}

func (r *PostgresClient) List(ctx context.Context) ([]domain.Slug, error) {
	rows, err := r.q.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list subscriptions: %w", err)
	}
	slugs := make([]domain.Slug, 0, len(rows))
	for _, row := range rows {
		slugs = append(slugs, domain.Slug{
			Owner: row.Owner,
			Repo:  row.Repo,
		})
	}
	return slugs, nil
}

func (r *PostgresClient) Push(ctx context.Context, slug domain.Slug) error {
	err := r.q.Push(ctx, db.PushParams{
		Owner: slug.Owner,
		Repo:  slug.Repo,
	})
	if err != nil {
		return fmt.Errorf("failed to push a sub: %w", err)
	}
	return nil
}

func (r *PostgresClient) Delete(ctx context.Context, slug domain.Slug) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin a transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	qtx := r.q.WithTx(tx)

	err = qtx.Delete(ctx, db.DeleteParams{
		Owner: slug.Owner,
		Repo:  slug.Repo,
	})
	if err != nil {
		return fmt.Errorf("failed to delete subscription %s/%s: %w", slug.Owner, slug.Repo, err)
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit changes: %w", err)
	}
	return nil
}

package postgres

import (
	"context"
	"fmt"
	"log/slog"

	db "repo-stat/processor/db/sqlc"
	"repo-stat/processor/internal/domain"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresClient struct {
	log  *slog.Logger
	pool *pgxpool.Pool
	q    *db.Queries
}

func NewPostgresClient(log *slog.Logger, ctx context.Context, dsn string) (*PostgresClient, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to init connection pool, %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	return &PostgresClient{
		log:  log,
		pool: pool,
		q:    db.New(pool),
	}, nil

}

func (r *PostgresClient) Close() {
	r.pool.Close()
}

func (r *PostgresClient) List(ctx context.Context) ([]domain.Repo, error) {
	rows, err := r.q.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list subscriptions: %w", err)
	}
	slugs := make([]domain.Repo, 0, len(rows))
	for _, row := range rows {
		r.log.Info("db adapter is called with list method",
			"values in the repo",
			fmt.Sprintf("Stars: %v, Forks: %v", row.Stars, row.Forks),
		)
		slugs = append(slugs, domain.Repo{
			Name:        row.Name,
			Description: row.Description.String,
			Stars:       int32(row.Stars.Int32),
			Forks:       int32(row.Forks.Int32),
			Date:        row.Date.Time,
		})
	}
	return slugs, nil
}

func (r *PostgresClient) Push(ctx context.Context, repo domain.Repo) error {
	var (
		desc  pgtype.Text
		stars pgtype.Int4
		forks pgtype.Int4
		date  pgtype.Timestamptz
	)
	r.log.Info("db adapter is called with push method",
		"values in the repo",
		fmt.Sprintf("Stars: %v, Forks: %v", repo.Stars, repo.Forks),
	)
	_ = desc.Scan(repo.Description)
	if err := stars.Scan(repo.Stars); err != nil {
		r.log.Error("scan failed", "error", err, "value", repo.Stars)
	}
	if err := forks.Scan(repo.Forks); err != nil {
		r.log.Error("scan failed", "error", err, "value", repo.Forks)
	}
	_ = date.Scan(repo.Date)
	r.log.Info("db adapter has called scan method to convert data types",
		"result of the scan",
		fmt.Sprintf("Stars: %v, Forks: %v", stars, forks),
	)
	err := r.q.Push(ctx, db.PushParams{
		Name:        repo.Name,
		Description: desc,
		Stars:       pgtype.Int4{Int32: int32(repo.Stars), Valid: true},
		Forks:       pgtype.Int4{Int32: int32(repo.Forks), Valid: true},
		Date:        date,
	})
	if err != nil {
		return fmt.Errorf("failed to push a repo: %w", err)
	}
	return nil
}

func (r *PostgresClient) Delete(ctx context.Context, repo domain.Repo) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin a transaction: %w", err)
	}
	defer func() {
		err := tx.Rollback(ctx)
		if err != nil {
			fmt.Printf("failed to begin a transaction: %v", err)
		}
	}()

	qtx := r.q.WithTx(tx)

	err = qtx.Delete(ctx, repo.Name)
	if err != nil {
		return fmt.Errorf("failed to delete information %s: %w", repo.Name, err)
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit changes: %w", err)
	}
	return nil
}

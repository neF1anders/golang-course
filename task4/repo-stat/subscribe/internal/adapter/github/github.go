package adapter

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"repo-stat/subscribe/internal/domain"
)

type Client struct {
	log  *slog.Logger
	conn *http.Client
}

func NewClient(log *slog.Logger) *Client {
	return &Client{
		conn: &http.Client{Timeout: 10 * time.Second},
		log:  log,
	}
}

func (c *Client) PingRepo(ctx context.Context, slug domain.Slug) (int, error) {
	url := fmt.Sprintf("https://github.com/%s/%s", slug.Owner, slug.Repo)
	resp, err := c.conn.Head(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	return resp.StatusCode, nil
}

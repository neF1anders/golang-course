package subscribe

import (
	"context"
	"log/slog"
	"repo-stat/api/internal/domain"

	subscribepb "repo-stat/proto/subscribe"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	log  *slog.Logger
	conn *grpc.ClientConn
	pb   subscribepb.SubscribeClient
}

func NewClient(address string, log *slog.Logger) (*Client, error) {
	conn, err := grpc.NewClient(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	return &Client{
		log:  log,
		conn: conn,
		pb:   subscribepb.NewSubscribeClient(conn),
	}, nil
}

func (c *Client) Subscribe(ctx context.Context, slug domain.Slug) error {
	_, err := c.pb.Subscribe(ctx, &subscribepb.Data{
		Owner: slug.Owner,
		Repo:  slug.Repo,
	})
	if err != nil {
		c.log.Error("subscribe sub failed", "error", err)
		return err
	}
	return nil
}
func (c *Client) Unsubscribe(ctx context.Context, slug domain.Slug) error {
	_, err := c.pb.Unsubscribe(ctx, &subscribepb.Data{
		Owner: slug.Owner,
		Repo:  slug.Repo,
	})
	if err != nil {
		c.log.Error("unsubscribe failed", "error", err)
		return err
	}
	return nil
}
func (c *Client) SubInfo(ctx context.Context) ([]domain.Slug, error) {
	data, err := c.pb.Subinfo(ctx, &subscribepb.Empty{})
	if err != nil {
		c.log.Error("subscribe info failed", "error", err)
		return nil, err
	}
	data_to_return := make([]domain.Slug, 0, len(data.Subscriptions))
	for _, el := range data.Subscriptions {
		data_to_return = append(data_to_return, domain.Slug{
			Owner: el.Owner,
			Repo:  el.Repo,
		})
	}
	return data_to_return, nil
}
func (c *Client) Close() error {
	return c.conn.Close()
}

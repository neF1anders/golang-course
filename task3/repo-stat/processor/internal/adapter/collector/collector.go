package collector

import (
	"context"
	"log/slog"

	"repo-stat/processor/internal/domain"
	collectorpb "repo-stat/proto/collector"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	log  *slog.Logger
	conn *grpc.ClientConn
	pb   collectorpb.CollectorClient
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
		pb:   collectorpb.NewCollectorClient(conn),
	}, nil
}

func (c *Client) GetInfo(ctx context.Context, owner, repo string) (domain.Repo, error) {
	data, err := c.pb.GetInfo(ctx, &collectorpb.Data{
		Owner: owner,
		Repo:  repo,
	})
	if err != nil {
		c.log.Error("collector fetch failed", "error", err)
		return domain.Repo{}, err
	}
	return domain.Repo{
		Name:        data.Name,
		Description: data.Description,
		Stars:       int(data.Stars),
		Forks:       int(data.Forks),
		Date:        data.Date,
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

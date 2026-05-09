package processor

import (
	"context"
	"log/slog"
	"repo-stat/api/internal/domain"

	processorpb "repo-stat/proto/processor"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	log  *slog.Logger
	conn *grpc.ClientConn
	pb   processorpb.ProcessorClient
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
		pb:   processorpb.NewProcessorClient(conn),
	}, nil
}

func (c *Client) Ping(ctx context.Context) (domain.PingStatus, error) {
	_, err := c.pb.Ping(ctx, &processorpb.PingRequest{})
	if err != nil {
		c.log.Error("processor ping failed", "error", err)
		return domain.PingStatusDown, err
	}

	return domain.PingStatusUp, nil
}

func (c *Client) GetInfo(ctx context.Context, owner, repo string) (domain.Repo, error) {
	data, err := c.pb.GetInfo(ctx, &processorpb.Data{
		Owner: owner,
		Repo:  repo,
	})
	if err != nil {
		c.log.Error("processor fetch failed", "error", err)
		return domain.Repo{}, err
	}
	return domain.Repo{
		Name:        data.Name,
		Description: data.Description,
		Stars:       int(data.Stars),
		Forks:       int(data.Forks),
		Date:        data.Date.AsTime(),
	}, nil
}

func (c *Client) GetSubInfo(ctx context.Context) ([]domain.Repo, error) {
	data, err := c.pb.GetSubInfo(ctx, &processorpb.Empty{})
	if err != nil {
		c.log.Error("processor sub-fetch failed", "error", err)
		return nil, err
	}
	res := make([]domain.Repo, 0, len(data.Subscriptions))
	for _, el := range data.Subscriptions {
		res = append(res, domain.Repo{
			Name:        el.Name,
			Description: el.Description,
			Stars:       int(el.Stars),
			Forks:       int(el.Forks),
			Date:        el.Date.AsTime(),
		})
	}
	return res, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

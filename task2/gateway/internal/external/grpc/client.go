package grpc

import (
	"context"
	"time"

	pb "gateway/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type CollectorClient struct {
	client pb.CollectorClient
	conn   *grpc.ClientConn
}

func NewCollectorClient(addr string) (*CollectorClient, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &CollectorClient{
		client: pb.NewCollectorClient(conn),
		conn:   conn,
	}, nil
}

func (c *CollectorClient) Close() error {
	return c.conn.Close()
}

func (c *CollectorClient) GetRepoInfo(owner, repo string) (*pb.Repo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return c.client.GetInfo(ctx, &pb.Data{
		Owner: owner,
		Repo:  repo,
	})
}

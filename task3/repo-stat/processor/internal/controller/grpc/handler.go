package grpc

import (
	"context"
	"log/slog"
	"repo-stat/processor/internal/usecase"
	pb "repo-stat/proto/processor"
)

type ProcessorServer struct {
	pb.UnimplementedProcessorServer
	log   *slog.Logger
	fetch *usecase.Fetch
	ping  *usecase.Ping
}

func NewProcessorServer(log *slog.Logger, fetch *usecase.Fetch, ping *usecase.Ping) *ProcessorServer {
	return &ProcessorServer{
		log:   log,
		fetch: fetch,
		ping:  ping,
	}
}

func (s *ProcessorServer) GetInfo(ctx context.Context, req *pb.Data) (*pb.Repo, error) {
	s.log.Debug("processor fetch request received")
	repoInfo, err := s.fetch.Execute(ctx, req.Owner, req.Repo)
	if err != nil {
		return nil, err
	}
	return &pb.Repo{
		Name:        repoInfo.Name,
		Description: repoInfo.Description,
		Stars:       int32(repoInfo.Stars),
		Forks:       int32(repoInfo.Forks),
		Date:        repoInfo.Date,
	}, nil
}

func (s *ProcessorServer) Ping(ctx context.Context, _ *pb.PingRequest) (*pb.PingResponse, error) {
	s.log.Debug("processor ping request received")
	return &pb.PingResponse{
		Reply: s.ping.Execute(ctx),
	}, nil
}

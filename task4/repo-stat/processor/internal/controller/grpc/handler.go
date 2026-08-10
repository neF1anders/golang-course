package grpc

import (
	"context"
	"log/slog"
	"repo-stat/processor/internal/usecase"
	pb "repo-stat/proto/processor"

	"google.golang.org/protobuf/types/known/timestamppb"
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
	repoInfo, err := s.fetch.GetInfo(ctx, req.Owner, req.Repo)
	if err != nil {
		return nil, err
	}
	return &pb.Repo{
		Name:        repoInfo.Name,
		Description: repoInfo.Description,
		Stars:       int32(repoInfo.Stars),
		Forks:       int32(repoInfo.Forks),
		Date:        timestamppb.New(repoInfo.Date),
	}, nil
}

func (s *ProcessorServer) GetSubInfo(ctx context.Context, req *pb.Empty) (*pb.Repos, error) {
	s.log.Debug("processor sub-fetch request received")
	repoInfo, err := s.fetch.GetSubInfo(ctx)
	if err != nil {
		return nil, err
	}
	res := &pb.Repos{}
	for _, el := range repoInfo {
		res.Subscriptions = append(res.Subscriptions, &pb.Repo{
			Name:        el.Name,
			Description: el.Description,
			Stars:       int32(el.Stars),
			Forks:       int32(el.Forks),
			Date:        timestamppb.New(el.Date),
		})
	}
	return res, nil
}

func (s *ProcessorServer) Ping(ctx context.Context, _ *pb.PingRequest) (*pb.PingResponse, error) {
	s.log.Debug("processor ping request received")
	return &pb.PingResponse{
		Reply: s.ping.Execute(ctx),
	}, nil
}

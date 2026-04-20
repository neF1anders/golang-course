package grpc

import (
	"context"
	"repo-stat/collector/internal/usecase"
	pb "repo-stat/proto/collector"
)

type CollectorServer struct {
	pb.UnimplementedCollectorServer
	getRepoInfo *usecase.GetRepoInfoUseCase
	getSubInfo  *usecase.GetSubInfoUseCase
}

func NewCollectorServer(getRepoInfo *usecase.GetRepoInfoUseCase, getSubInfo *usecase.GetSubInfoUseCase) *CollectorServer {
	return &CollectorServer{
		getRepoInfo: getRepoInfo,
		getSubInfo:  getSubInfo,
	}
}

func (s *CollectorServer) GetInfo(ctx context.Context, req *pb.Data) (*pb.Repo, error) {
	repoInfo, err := s.getRepoInfo.Execute(req.Owner, req.Repo)
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

func (s *CollectorServer) GetSubInfo(ctx context.Context, req *pb.Empty) (*pb.Repos, error) {
	repoInfo, err := s.getSubInfo.Execute(ctx)
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
			Date:        el.Date,
		})
	}
	return res, nil
}

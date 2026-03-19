package grpc

import (
	"collector/internal/usecase"
	pb "collector/proto"
	"context"
)

type CollectorServer struct {
	pb.UnimplementedCollectorServer
	getRepoInfo *usecase.GetRepoInfoUseCase
}

func NewCollectorServer(getRepoInfo *usecase.GetRepoInfoUseCase) *CollectorServer {
	return &CollectorServer{
		getRepoInfo: getRepoInfo,
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

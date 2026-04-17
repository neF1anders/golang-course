package grpc

import (
	"context"
	"log/slog"
	pb "repo-stat/proto/subscribe"
	"repo-stat/subscribe/internal/domain"
	"repo-stat/subscribe/internal/usecase"
)

type SubscribeServer struct {
	pb.UnimplementedSubscribeServer
	log  *slog.Logger
	db   *usecase.DBUseCase
	ping *usecase.PingUseCase
}

func NewSubscribeServer(log *slog.Logger, db *usecase.DBUseCase, ping *usecase.PingUseCase) *SubscribeServer {
	return &SubscribeServer{
		log:  log,
		db:   db,
		ping: ping,
	}
}

func (s *SubscribeServer) Subinfo(ctx context.Context, empty *pb.Empty) (*pb.InfoReply, error) {
	s.log.Debug("subscribe info request received")
	info, err := s.db.List(ctx)
	if err != nil {
		return nil, err
	}
	data := make([]*pb.Data, len(info))
	for _, el := range info {
		data = append(data, &pb.Data{
			Owner: el.Owner,
			Repo:  el.Owner,
		})
	}
	return &pb.InfoReply{
		Subscriptions: data,
	}, nil
}
func (s *SubscribeServer) Subscribe(ctx context.Context, slug *pb.Data) (*pb.SubReply, error) {
	s.log.Debug("subscribe request received")
	err := s.db.Push(ctx, domain.Slug{
		Owner: slug.Owner,
		Repo:  slug.Repo,
	})
	if err != nil {
		return &pb.SubReply{
			Reply: err.Error(),
		}, err
	}
	return &pb.SubReply{
		Reply: "OK",
	}, nil
}
func (s *SubscribeServer) Unsubscribe(ctx context.Context, slug *pb.Data) (*pb.SubReply, error) {
	s.log.Debug("unsubscribe request received")
	err := s.db.Delete(ctx, domain.Slug{
		Owner: slug.Owner,
		Repo:  slug.Repo,
	})
	if err != nil {
		return &pb.SubReply{
			Reply: err.Error(),
		}, err
	}
	return &pb.SubReply{
		Reply: "OK",
	}, nil
}

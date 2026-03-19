package main

import (
	"log"
	"net"

	"google.golang.org/grpc"

	deliverygrpc "collector/internal/delivery/grpc"
	"collector/internal/external/github"
	"collector/internal/usecase"
	pb "collector/proto"
)

func main() {
	githubClient := github.NewClient()
	getRepoInfo := usecase.NewGetRepoInfoUseCase(githubClient)
	server := grpc.NewServer()
	collectorServer := deliverygrpc.NewCollectorServer(getRepoInfo)
	pb.RegisterCollectorServer(server, collectorServer)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Println("Collector gRPC server listening on :50051")
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}

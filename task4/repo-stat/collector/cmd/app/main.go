package main

import (
	"log"
	"net"

	"google.golang.org/grpc"

	"repo-stat/collector/internal/adapter/github"
	deliverygrpc "repo-stat/collector/internal/controller/grpc"
	"repo-stat/collector/internal/usecase"
	pb "repo-stat/proto/collector"
)

func main() { //config unused - i want to, but have some difficulties
	githubClient := github.NewClient()
	getRepoInfo := usecase.NewGetRepoInfoUseCase(githubClient)
	server := grpc.NewServer()
	collectorServer := deliverygrpc.NewCollectorServer(getRepoInfo)
	pb.RegisterCollectorServer(server, collectorServer)

	lis, err := net.Listen("tcp", ":8084")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Println("Collector gRPC server listening on :8084")
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}

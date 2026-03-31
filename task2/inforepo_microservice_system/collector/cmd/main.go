package main

import (
	"collector/internal/adapter/github"
	grpchandler "collector/internal/delivery/grpc"
	"collector/internal/usecase"
	pb "collector/pb"
	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {
	githubAdapter := github.NewAdapter()
	service := usecase.NewService(githubAdapter)
	grpcHandler := grpchandler.NewHandler(service)

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer()
	pb.RegisterCollectorServiceServer(server, grpcHandler)

	log.Println("Collector Service listening on port 50051")
	if err := server.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

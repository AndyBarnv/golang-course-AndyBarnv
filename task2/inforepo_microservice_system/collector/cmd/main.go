package main

import (
	"collector/internal/adapter/github"
	"collector/internal/delivery/grpc"
	"collector/internal/usecase"
	"fmt"
	pb "inforepo_microservice_system/proto/pb"
	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {
	githubAdapter := github.NewAdapter()
	service := usecase.NewService(githubAdapter)
	grpcHandler := grpc.NewHandler(service)

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer()
	pb.RegisterCollectorServiceServer(server, grpcHandler)

	fmt.Println("Collector Service listening on port 50051") // Потом убраттть
	if err := server.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

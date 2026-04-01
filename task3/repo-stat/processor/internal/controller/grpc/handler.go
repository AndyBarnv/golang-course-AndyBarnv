package grpc

import (
	"context"
	"repo-stat/processor/internal/usecase"
	processorpblib "repo-stat/proto/processor"
)

type Server struct {
	processorpblib.UnimplementedProcessorServiceServer
	uc *usecase.Service
}

func NewServer(uc *usecase.Service) *Server {
	return &Server{uc: uc}
}

func (s *Server) GetRepositoryInfo(ctx context.Context, req *processorpblib.RepositoryRequest) (*processorpblib.RepositoryResponse, error) {
	return s.uc.GetRepoInfo(ctx, req.Owner, req.Repo)
}

func (s *Server) Ping(ctx context.Context, req *processorpblib.PingRequest) (*processorpblib.PingResponse, error) {
	return &processorpblib.PingResponse{}, nil
}

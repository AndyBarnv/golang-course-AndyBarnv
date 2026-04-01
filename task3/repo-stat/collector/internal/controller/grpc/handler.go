package grpc

import (
	"context"
	"repo-stat/collector/internal/usecase"
	collectorpblib "repo-stat/proto/collector"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	collectorpblib.UnimplementedCollectorServiceServer
	uc *usecase.Service
}

func NewServer(uc *usecase.Service) *Server {
	return &Server{uc: uc}
}

func (s *Server) GetRepositoryInfo(ctx context.Context, req *collectorpblib.RepositoryRequest) (*collectorpblib.RepositoryResponse, error) {
	repo, err := s.uc.GetRepo(ctx, req.Owner, req.Repo)
	if err != nil {
		if err.Error() == "repository not found" {
			return nil, status.Errorf(codes.NotFound, err.Error())
		}
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &collectorpblib.RepositoryResponse{
		FullName:    repo.FullName,
		Description: repo.Description,
		Stars:       repo.Stars,
		Forks:       repo.Forks,
		CreatedAt:   repo.CreatedAt,
	}, nil
}

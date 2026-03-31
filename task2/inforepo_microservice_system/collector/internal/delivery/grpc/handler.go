package grpc

import (
	"collector/internal/usecase"
	"collector/pb"
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handler struct {
	pb.UnimplementedCollectorServiceServer
	usecase *usecase.Service
}

func NewHandler(uc *usecase.Service) *Handler {
	return &Handler{usecase: uc}
}

func (h *Handler) GetRepositoryInfo(ctx context.Context, req *pb.RepositoryRequest) (*pb.RepositoryResponse, error) {
	repo, err := h.usecase.GetRepoInfo(ctx, req.Owner, req.Repo)
	if err != nil {
		if err.Error() == "repository not found" {
			return nil, status.Errorf(codes.NotFound, err.Error())
		}
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &pb.RepositoryResponse{
		Name:        repo.Name,
		Description: repo.Description,
		Stars:       int32(repo.Stars),
		Forks:       int32(repo.Forks),
		CreatedAt:   repo.CreationDate.Format(time.RFC3339),
	}, nil
}

package usecase

import (
	"context"
	"gateway/internal/domain"
)

type CollectorClient interface {
	GetRepo(ctx context.Context, owner, repo string) (*domain.Repository, error)
}

type Service struct {
	collector CollectorClient
}

func (s *Service) FetchRepository(ctx context.Context, owner, repo string) (*domain.Repository, error) {
	return s.collector.GetRepo(ctx, owner, repo)
}

func NewService(client CollectorClient) *Service {
	return &Service{collector: client}
}

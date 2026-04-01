package usecase

import (
	"context"
	"repo-stat/collector/internal/domain"
)

type RepoGetter interface {
	GetRepo(ctx context.Context, owner, repo string) (*domain.Repository, error)
}

type Service struct {
	getter RepoGetter
}

func NewService(getter RepoGetter) *Service {
	return &Service{getter: getter}
}

func (s *Service) GetRepo(ctx context.Context, owner, repo string) (*domain.Repository, error) {
	return s.getter.GetRepo(ctx, owner, repo)
}

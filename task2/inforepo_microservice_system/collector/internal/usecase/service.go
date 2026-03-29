package usecase

import (
	"collector/internal/domain"
	"context"
)

type GithubAdapter interface {
	FetchRepo(ctx context.Context, owner, repo string) (*domain.Repository, error)
}

type Service struct {
	githubAdapter GithubAdapter
}

func (s *Service) GetRepoInfo(ctx context.Context, owner, repo string) (*domain.Repository, error) {
	return s.githubAdapter.FetchRepo(ctx, owner, repo)
}

func NewService(adapter GithubAdapter) *Service {
	return &Service{githubAdapter: adapter}
}

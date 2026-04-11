package usecase

import (
	"context"
	processorpblib "repo-stat/proto/processor"
)

type RepoGetter interface {
	GetRepoInfo(ctx context.Context, owner, repo string) (*processorpblib.RepositoryResponse, error)
}

type Service struct {
	getter RepoGetter
}

func NewService(getter RepoGetter) *Service {
	return &Service{getter: getter}
}

func (s *Service) GetRepoInfo(ctx context.Context, owner, repo string) (*processorpblib.RepositoryResponse, error) {
	return s.getter.GetRepoInfo(ctx, owner, repo)
}

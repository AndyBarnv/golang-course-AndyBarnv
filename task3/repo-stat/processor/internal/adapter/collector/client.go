package collector

import (
	"context"
	collectorpblib "repo-stat/proto/collector"
	processorpblib "repo-stat/proto/processor"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	pb collectorpblib.CollectorServiceClient
}

func NewClient(address string) (*Client, error) {
	conn, err := grpc.NewClient(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	return &Client{
		pb: collectorpblib.NewCollectorServiceClient(conn),
	}, nil
}

func (c *Client) GetRepoInfo(ctx context.Context, owner, repo string) (*processorpblib.RepositoryResponse, error) {
	resp, err := c.pb.GetRepositoryInfo(ctx, &collectorpblib.RepositoryRequest{
		Owner: owner,
		Repo:  repo,
	})
	if err != nil {
		return nil, err
	}

	return &processorpblib.RepositoryResponse{
		FullName:    resp.FullName,
		Description: resp.Description,
		Stars:       resp.Stars,
		Forks:       resp.Forks,
		CreatedAt:   resp.CreatedAt,
	}, nil
}

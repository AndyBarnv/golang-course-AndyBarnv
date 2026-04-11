package processor

import (
	"context"
	"repo-stat/proto/processor"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	pb processor.ProcessorServiceClient
}

func NewClient(address string) (*Client, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &Client{pb: processor.NewProcessorServiceClient(conn)}, nil
}

func (c *Client) Ping(ctx context.Context) error {
	_, err := c.pb.Ping(ctx, &processor.PingRequest{})
	return err
}

func (c *Client) GetRepoInfo(ctx context.Context, owner, repo string) (*processor.RepositoryResponse, error) {
	return c.pb.GetRepositoryInfo(ctx, &processor.RepositoryRequest{Owner: owner, Repo: repo})
}

package grpc

import (
	"context"
	"gateway/internal/domain"
	"gateway/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Adapter struct {
	client pb.CollectorServiceClient
	conn   *grpc.ClientConn
}

func NewAdapter(collectorAddr string) (*Adapter, error) {
	conn, err := grpc.Dial(collectorAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	client := pb.NewCollectorServiceClient(conn)
	return &Adapter{client: client, conn: conn}, nil
}

func (a *Adapter) Close() {
	a.conn.Close()
}

func (a *Adapter) GetRepo(ctx context.Context, owner, repo string) (*domain.Repository, error) {
	resp, err := a.client.GetRepositoryInfo(ctx, &pb.RepositoryRequest{
		Owner: owner,
		Repo:  repo,
	})
	if err != nil {
		return nil, err
	}

	return &domain.Repository{
		Name:        resp.Name,
		Description: resp.Description,
		Stars:       int(resp.Stars),
		Forks:       int(resp.Forks),
		CreatedAt:   resp.CreatedAt,
	}, nil
}

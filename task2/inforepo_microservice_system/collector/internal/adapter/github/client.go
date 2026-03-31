package github

import (
	"collector/internal/domain"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"time"
)

type Adapter struct {
	client *http.Client
}

var githubResp struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Stargazers  int       `json:"stargazers_count"`
	Forks       int       `json:"forks_count"`
	CreatedAt   time.Time `json:"created_at"`
}

func NewAdapter() *Adapter {
	return &Adapter{
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

func (a *Adapter) FetchRepo(ctx context.Context, owner, repo string) (*domain.Repository, error) {
	baseURL, _ := url.Parse("https://api.github.com")
	baseURL.Path = path.Join("repos", owner, repo)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Go-Microservice-Collector")

	resp, err := a.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, errors.New("repository not found")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("github api error: status %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&githubResp); err != nil {
		return nil, err
	}

	return &domain.Repository{
		Name:         githubResp.Name,
		Description:  githubResp.Description,
		Stars:        githubResp.Stargazers,
		Forks:        githubResp.Forks,
		CreationDate: githubResp.CreatedAt,
	}, nil
}

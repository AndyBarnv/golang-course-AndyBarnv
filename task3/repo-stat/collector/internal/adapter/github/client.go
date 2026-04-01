package github

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"repo-stat/collector/internal/domain"
	"time"
)

type Client struct {
	httpClient *http.Client
}

func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *Client) GetRepo(ctx context.Context, owner, repo string) (*domain.Repository, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s", owner, repo)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("repository not found")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("github api error: %d", resp.StatusCode)
	}

	var ghResp struct {
		FullName    string `json:"full_name"`
		Description string `json:"description"`
		Stars       int64  `json:"stargazers_count"`
		Forks       int64  `json:"forks_count"`
		CreatedAt   string `json:"created_at"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&ghResp); err != nil {
		return nil, err
	}

	return &domain.Repository{
		FullName:    ghResp.FullName,
		Description: ghResp.Description,
		Stars:       ghResp.Stars,
		Forks:       ghResp.Forks,
		CreatedAt:   ghResp.CreatedAt,
	}, nil
}

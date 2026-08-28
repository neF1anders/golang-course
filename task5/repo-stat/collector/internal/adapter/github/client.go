package github

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	entity "repo-stat/collector/internal/domain"
)

type Client struct {
	log        *slog.Logger
	httpClient *http.Client
}

func NewClient(log *slog.Logger) *Client {
	return &Client{
		log:        log,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

type githubResponse struct {
	Name        string    `json:"full_name"`
	Description string    `json:"description"`
	Stars       int       `json:"stargazers_count"`
	Forks       int       `json:"forks_count"`
	CreatedAt   time.Time `json:"created_at"`
}

func (c *Client) Fetch(owner, repo string) (*entity.Repo, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s", owner, repo)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("User-Agent", "repo-stat-collector")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			c.log.Error("error during closing responce body in fetcher", "error", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		c.log.Debug("API return non-OK status", "debug", resp.StatusCode)
		return nil, fmt.Errorf("github api return: %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.log.Error("error reading response body", "error", err)
		return nil, fmt.Errorf("reading response error: %d", err)
	}
	var response githubResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		c.log.Error("error retrieving from JSON", "error", err)
		return nil, fmt.Errorf("retriever returned: %d", err)
	}
	return &entity.Repo{
		Name:        response.Name,
		Description: response.Description,
		Stars:       response.Stars,
		Forks:       response.Forks,
		Date:        response.CreatedAt,
	}, nil
}

package github

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"repo-stat/collector/internal/entity"
)

type Client struct {
	httpClient *http.Client
}

func NewClient() *Client {
	return &Client{
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

	resp, err := c.httpClient.Get(url)

	if err != nil {
		return nil, err
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Printf("Error during closing responce body in fetcher: %v", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		log.Printf("API return non-OK status: %v", resp.Status)
		return nil, fmt.Errorf("github api return: %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return nil, fmt.Errorf("reading response error: %d", err)
	}
	var response githubResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("Error retrieving from JSON: %v", err)
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

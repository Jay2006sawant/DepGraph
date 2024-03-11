package github

import (
	"context"
	"fmt"
	"os"

	"github.com/google/go-github/v45/github"
	"golang.org/x/oauth2"
)

// Client wraps the GitHub API client
type Client struct {
	client *github.Client
	ctx    context.Context
}

// NewClient creates a new GitHub API client using the token from environment
func NewClient() (*Client, error) {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("GITHUB_TOKEN environment variable is not set")
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	return &Client{
		client: client,
		ctx:    ctx,
	}, nil
}

// ListRepositories returns all repositories for the given organization
func (c *Client) ListRepositories(org string) ([]*github.Repository, error) {
	opt := &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}

	var allRepos []*github.Repository
	for {
		repos, resp, err := c.client.Repositories.ListByOrg(c.ctx, org, opt)
		if err != nil {
			return nil, fmt.Errorf("failed to list repositories: %v", err)
		}
		allRepos = append(allRepos, repos...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	return allRepos, nil
}

// GetModuleFile fetches the go.mod file content from a repository
func (c *Client) GetModuleFile(owner, repo, path string) (string, error) {
	content, _, _, err := c.client.Repositories.GetContents(
		c.ctx,
		owner,
		repo,
		path,
		nil,
	)
	if err != nil {
		return "", fmt.Errorf("failed to get go.mod: %v", err)
	}

	fileContent, err := content.GetContent()
	if err != nil {
		return "", fmt.Errorf("failed to decode content: %v", err)
	}

	return fileContent, nil
} 
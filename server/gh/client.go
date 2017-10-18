package gh

import (
	"context"

	"github.com/google/go-github/github"
)

// Client refers to a client which wishes to connect to specific repository of a user.
type Client struct {
	*github.Client
	owner string
	repo  string
}

// NewClient constructs a new instance of Client
func NewClient(owner, repo string) *Client {
	return &Client{
		Client: github.NewClient(nil),
		owner:  owner,
		repo:   repo,
	}
}

// GetIssues gets issues of a repo.
func (c *Client) GetIssues(ctx context.Context) ([]*github.Issue, error) {
	issues, _, err := c.Client.Issues.ListByRepo(ctx, c.owner, c.repo, nil)
	if err != nil {
		return nil, err
	}
	return issues, nil
}

// GetPullRequests gets pull request list for a repo.
func (c *Client) GetPullRequests(ctx context.Context) ([]*github.PullRequest, error) {
	pullRequests, _, err := c.Client.PullRequests.List(ctx, c.owner, c.repo, nil)
	if err != nil {
		return nil, err
	}
	return pullRequests, nil
}

// GetAllLabels gets all labels of a repo, not an issue, nor a pull request
func (c *Client) GetAllLabels(ctx context.Context) ([]*github.Label, error) {
	labels, _, err := c.Client.Issues.ListLabels(ctx, c.owner, c.repo, nil)
	if err != nil {
		return nil, err
	}
	return labels, nil
}

// GetLabelsByIssue gets labels attached on a single issue whose id is num.
func (c *Client) GetLabelsByIssue(ctx context.Context, num int) ([]*github.Label, error) {
	labels, _, err := c.Client.Issues.ListLabelsByIssue(ctx, c.owner, c.repo, num, nil)
	if err != nil {
		return nil, err
	}
	return labels, nil
}

// AddLabelsToIssue adds labels to an issue
func (c *Client) AddLabelsToIssue(ctx context.Context, num int, labels []string) error {
	if _, _, err := c.Client.Issues.AddLabelsToIssue(ctx, c.owner, c.repo, num, labels); err != nil {
		return err
	}
	return nil
}

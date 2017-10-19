package gh

import (
	"context"
	"sync"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// Client refers to a client which wishes to connect to specific repository of a user.
type Client struct {
	sync.Mutex
	*github.Client
	owner string
	repo  string
}

// NewClient constructs a new instance of Client
func NewClient(owner, repo, token string) *Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{
			AccessToken: token,
		},
	)
	tc := oauth2.NewClient(ctx, ts)

	return &Client{
		Client: github.NewClient(tc),
		owner:  owner,
		repo:   repo,
	}
}

// GetIssues gets issues of a repo.
func (c *Client) GetIssues(ctx context.Context) ([]*github.Issue, error) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	issues, _, err := c.Client.Issues.ListByRepo(ctx, c.owner, c.repo, nil)
	if err != nil {
		return nil, err
	}
	return issues, nil
}

// GetPullRequests gets pull request list for a repo.
func (c *Client) GetPullRequests(ctx context.Context) ([]*github.PullRequest, error) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	pullRequests, _, err := c.Client.PullRequests.List(ctx, c.owner, c.repo, nil)
	if err != nil {
		return nil, err
	}
	return pullRequests, nil
}

// GetAllLabels gets all labels of a repo, not an issue, nor a pull request
func (c *Client) GetAllLabels(ctx context.Context) ([]*github.Label, error) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	labels, _, err := c.Client.Issues.ListLabels(ctx, c.owner, c.repo, nil)
	if err != nil {
		return nil, err
	}
	return labels, nil
}

// GetLabelsByIssue gets labels attached on a single issue whose id is num.
func (c *Client) GetLabelsByIssue(ctx context.Context, num int) ([]*github.Label, error) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	labels, _, err := c.Client.Issues.ListLabelsByIssue(ctx, c.owner, c.repo, num, nil)
	if err != nil {
		return nil, err
	}
	return labels, nil
}

// AddLabelsToIssue adds labels to an issue
func (c *Client) AddLabelsToIssue(ctx context.Context, num int, labels []string) error {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	if _, _, err := c.Client.Issues.AddLabelsToIssue(ctx, c.owner, c.repo, num, labels); err != nil {
		return err
	}
	return nil
}

// AddCommentToIssue adds comment to an issue
func (c *Client) AddCommentToIssue(ctx context.Context, num int, comment *github.IssueComment) error {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	if _, _, err := c.Client.Issues.CreateComment(ctx, c.owner, c.repo, num, comment); err != nil {
		return err
	}
	return nil
}

// AssignIssueToUsers assigns users to the specified issue.
func (c *Client) AssignIssueToUsers(ctx context.Context, num int, users []string) error {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	if _, _, err := c.Client.Issues.AddAssignees(ctx, c.owner, c.repo, num, users); err != nil {
		return err
	}
	return nil
}

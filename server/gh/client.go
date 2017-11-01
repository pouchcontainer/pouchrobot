package gh

import (
	"context"
	"sync"

	"github.com/sirupsen/logrus"

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
func (c *Client) GetIssues(ctx context.Context, opt *github.IssueListByRepoOptions) ([]*github.Issue, error) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	issues, _, err := c.Client.Issues.ListByRepo(ctx, c.owner, c.repo, opt)
	if err != nil {
		logrus.Errorf("failed to list issues in repo %s: %v", c.repo, err)
		return nil, err
	}
	logrus.Debugf("succeed in getting issues in repo %s", c.repo)
	return issues, nil
}

// GetPullRequests gets pull request list for a repo.
func (c *Client) GetPullRequests(ctx context.Context, opt *github.PullRequestListOptions) ([]*github.PullRequest, error) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	pullRequests, _, err := c.Client.PullRequests.List(ctx, c.owner, c.repo, opt)
	if err != nil {
		logrus.Errorf("failed to list pull request in repo %s: %v", c.repo, err)
		return nil, err
	}
	logrus.Debugf("succeed in getting pull requests in repo %s", c.repo)
	return pullRequests, nil
}

// GetAllLabels gets all labels of a repo, not an issue, nor a pull request
func (c *Client) GetAllLabels(ctx context.Context) ([]*github.Label, error) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	labels, _, err := c.Client.Issues.ListLabels(ctx, c.owner, c.repo, nil)
	if err != nil {
		logrus.Errorf("failed to get labels in repo %s: %v", c.repo, err)
		return nil, err
	}
	logrus.Debugf("succeed in listing all labels in repo %s", c.repo)
	return labels, nil
}

// GetLabelsInIssue gets labels attached on a single issue whose id is num.
func (c *Client) GetLabelsInIssue(ctx context.Context, num int) ([]*github.Label, error) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	labels, _, err := c.Client.Issues.ListLabelsByIssue(ctx, c.owner, c.repo, num, nil)
	if err != nil {
		logrus.Errorf("failed to get labels in issue %d: %v", num, err)
		return nil, err
	}
	logrus.Debugf("succeed in getting labels in issue %d", num)
	return labels, nil
}

// AddLabelsToIssue adds labels to an issue
func (c *Client) AddLabelsToIssue(ctx context.Context, num int, labels []string) error {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	if _, _, err := c.Client.Issues.AddLabelsToIssue(ctx, c.owner, c.repo, num, labels); err != nil {
		logrus.Errorf("failed to add labels %s to issue(pr) %d: %v", labels, num, err)
		return err
	}
	logrus.Debugf("succeed in adding labels %v for issue %d", labels, num)
	return nil
}

// RemoveLabelForIssue removes a label from an issue
func (c *Client) RemoveLabelForIssue(ctx context.Context, num int, label string) error {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	if _, err := c.Client.Issues.RemoveLabelForIssue(ctx, c.owner, c.repo, num, label); err != nil {
		logrus.Errorf("failed to remove label %s for issue(pr) %d: %v", label, num, err)
		return err
	}
	logrus.Debugf("succeed in removing label %v for issue %d", label, num)
	return nil
}

// ListPRComments lists comments for a pull request
func (c *Client) ListPRComments(ctx context.Context, num int) ([]*github.PullRequestComment, error) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	prComments, _, err := c.Client.PullRequests.ListComments(ctx, c.owner, c.repo, num, nil)
	if err != nil {
		logrus.Errorf("failed to list comments for pr %d: %v", num, err)
		return nil, err
	}
	logrus.Debugf("succeed in list comments for pr %d:", num)
	return prComments, nil
}

// AddCommentToIssue adds comment to an issue
func (c *Client) AddCommentToIssue(ctx context.Context, num int, comment *github.IssueComment) error {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	if _, _, err := c.Client.Issues.CreateComment(ctx, c.owner, c.repo, num, comment); err != nil {
		logrus.Errorf("failed to add comment %s to issue(pr) %d: %v", *(comment.Body), num, err)
		return err
	}
	logrus.Debugf("succeed in adding comment %s for issue %d", *(comment.Body), num)
	return nil
}

// AddCommentToPR adds comment to a pull request
func (c *Client) AddCommentToPR(ctx context.Context, num int, comment *github.IssueComment) error {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	if _, _, err := c.Client.Issues.CreateComment(ctx, c.owner, c.repo, num, comment); err != nil {
		logrus.Errorf("failed to add comment %s to pr %d: %v", *(comment.Body), num, err)
		return err
	}
	logrus.Debugf("succeed in creating comment %s for pull request %d", *(comment.Body), num)
	return nil
}

// RemoveCommentForPR removes a comment for a pull request
func (c *Client) RemoveCommentForPR(ctx context.Context, num int) error {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	if _, err := c.Client.PullRequests.DeleteComment(ctx, c.owner, c.repo, num); err != nil {
		logrus.Errorf("failed to remove comment %d: %v", num, err)
		return err
	}
	logrus.Debugf("succeed in removing comment %s for pull request", num)
	return nil
}

// AssignIssueToUsers assigns users to the specified issue.
func (c *Client) AssignIssueToUsers(ctx context.Context, num int, users []string) error {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	if _, _, err := c.Client.Issues.AddAssignees(ctx, c.owner, c.repo, num, users); err != nil {
		logrus.Errorf("failed to assign users %s to issue(pr) %d: %v", users, num, err)
		return err
	}
	logrus.Debugf("succeed in assign users %s for pull request %d", users, num)
	return nil
}

// UnassignIssueToUsers assigns users to the specified issue.
func (c *Client) UnassignIssueToUsers(ctx context.Context, num int, users []string) error {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	if _, _, err := c.Client.Issues.AddAssignees(ctx, c.owner, c.repo, num, users); err != nil {
		logrus.Errorf("failed to assign users %s to issue(pr) %d: %v", users, num, err)
		return err
	}
	logrus.Debugf("succeed in assign users %s for pull request %d", users, num)
	return nil
}

// IssueHasLabel judges if an issue has a specified label.
func (c *Client) IssueHasLabel(num int, inputLabel string) bool {
	labels, err := c.GetLabelsInIssue(context.Background(), num)
	if err != nil {
		return false
	}
	for _, label := range labels {
		if label.GetName() == inputLabel {
			return true
		}
	}
	return false
}

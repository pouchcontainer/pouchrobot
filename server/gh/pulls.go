package gh

import (
	"context"

	"github.com/google/go-github/github"
	"github.com/sirupsen/logrus"
)

// GetPullRequests gets pull request list for a repo.
func (c *Client) GetPullRequests(opt *github.PullRequestListOptions) ([]*github.PullRequest, error) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	pullRequests, _, err := c.Client.PullRequests.List(context.Background(), c.owner, c.repo, opt)
	if err != nil {
		logrus.Errorf("failed to list pull request in repo %s: %v", c.repo, err)
		return nil, err
	}
	logrus.Debugf("succeed in getting pull requests in repo %s", c.repo)
	return pullRequests, nil
}

// GetSinglePR gets a single PR from repo.
func (c *Client) GetSinglePR(num int) (*github.PullRequest, error) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	pullRequest, _, err := c.Client.PullRequests.Get(context.Background(), c.owner, c.repo, num)
	if err != nil {
		logrus.Errorf("failed to get single pull request %d in repo %s: %v", num, c.repo, err)
		return nil, err
	}
	logrus.Debugf("succeed in getting single pull request %d in repo %s", num, c.repo)
	return pullRequest, nil
}

// ListPRComments lists comments for a pull request.
func (c *Client) ListPRComments(num int) ([]*github.PullRequestComment, error) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	prComments, _, err := c.Client.PullRequests.ListComments(context.Background(), c.owner, c.repo, num, nil)
	if err != nil {
		logrus.Errorf("failed to list comments for pr %d: %v", num, err)
		return nil, err
	}
	logrus.Debugf("succeed in list comments for pr %d:", num)
	return prComments, nil
}

// AddCommentToPR adds comment to a pull request.
func (c *Client) AddCommentToPR(num int, comment *github.IssueComment) error {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	if _, _, err := c.Client.Issues.CreateComment(context.Background(), c.owner, c.repo, num, comment); err != nil {
		logrus.Errorf("failed to add comment %s to pr %d: %v", *(comment.Body), num, err)
		return err
	}
	logrus.Debugf("succeed in creating comment %s for pull request %d", *(comment.Body), num)
	return nil
}

// ListCommits lists all commits in a pull request.
func (c *Client) ListCommits(num int) ([]*github.RepositoryCommit, error) {
	commits, _, err := c.PullRequests.ListCommits(context.Background(), c.owner, c.repo, num, nil)
	if err != nil {
		logrus.Errorf("failed to list commits in pull request %d: %v", num, err)
		return nil, err
	}
	logrus.Debugf("succeed in listing commits in pull request %d", num)
	return commits, nil
}

// RemoveCommentForPR removes a comment for a pull request.
func (c *Client) RemoveCommentForPR(num int) error {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	if _, err := c.Client.PullRequests.DeleteComment(context.Background(), c.owner, c.repo, num); err != nil {
		logrus.Errorf("failed to remove comment %d: %v", num, err)
		return err
	}
	logrus.Debugf("succeed in removing comment %s for pull request", num)
	return nil
}

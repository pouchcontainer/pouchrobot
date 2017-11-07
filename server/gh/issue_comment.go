package gh

import (
	"context"
	"strings"

	"github.com/google/go-github/github"
	"github.com/sirupsen/logrus"
)

// ListComments lists all comments in an issue including pull request.
func (c *Client) ListComments(num int) ([]*github.IssueComment, error) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	comments, _, err := c.Client.Issues.ListComments(context.Background(), c.owner, c.repo, num, nil)
	if err != nil {
		logrus.Errorf("failed to list comment in issue(pr) %d: %v", num, err)
		return nil, err
	}
	logrus.Debugf("succeed in listing comments for issue(pr) %d", num)
	return comments, nil
}

// AddCommentToIssue adds comment to an issue.
func (c *Client) AddCommentToIssue(num int, comment *github.IssueComment) error {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	if _, _, err := c.Client.Issues.CreateComment(context.Background(), c.owner, c.repo, num, comment); err != nil {
		logrus.Errorf("failed to add comment %s to issue(pr) %d: %v", *(comment.Body), num, err)
		return err
	}
	logrus.Debugf("succeed in adding comment %s for issue %d", *(comment.Body), num)
	return nil
}

// RemoveComment removes a comment for an issue.
func (c *Client) RemoveComment(id int) error {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	if _, err := c.Client.Issues.DeleteComment(context.Background(), c.owner, c.repo, id); err != nil {
		logrus.Errorf("failed to remove comment %d: %v", id, err)
		return err
	}
	logrus.Debugf("succeed in removing comment %s for issue", id)
	return nil
}

// RemoveCommentViaString removes a comment in an issue via given string
func (c *Client) RemoveCommentViaString(num int, str string) error {
	comments, err := c.ListComments(num)
	if err != nil {
		return err
	}

	for _, comment := range comments {
		if comment.Body != nil && strings.Contains(*(comment.Body), str) {
			return c.RemoveComment(*(comment.ID))
		}
	}
	return nil
}

// IssueHasComment returns true if the issue contains a commnet who has substring of 'elment'
func (c *Client) IssueHasComment(num int, element string) (int, bool) {
	comments, err := c.ListComments(num)
	if err != nil {
		return -1, false
	}
	for _, comment := range comments {
		if comment.Body != nil && strings.Contains(*(comment.Body), element) {
			return *(comment.ID), true
		}
	}
	return -1, false
}

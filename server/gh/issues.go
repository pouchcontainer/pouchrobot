package gh

import (
	"context"

	"github.com/google/go-github/github"
	"github.com/sirupsen/logrus"
)

// GetIssues gets issues of a repo.
func (c *Client) GetIssues(opt *github.IssueListByRepoOptions) ([]*github.Issue, error) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	issues, _, err := c.Client.Issues.ListByRepo(context.Background(), c.owner, c.repo, opt)
	if err != nil {
		logrus.Errorf("failed to list issues in repo %s: %v", c.repo, err)
		return nil, err
	}
	logrus.Debugf("succeed in getting issues in repo %s", c.repo)
	return issues, nil
}

// GetAllLabels gets all labels of a repo, not an issue, nor a pull request
func (c *Client) GetAllLabels() ([]*github.Label, error) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	labels, _, err := c.Client.Issues.ListLabels(context.Background(), c.owner, c.repo, nil)
	if err != nil {
		logrus.Errorf("failed to get labels in repo %s: %v", c.repo, err)
		return nil, err
	}
	logrus.Debugf("succeed in listing all labels in repo %s", c.repo)
	return labels, nil
}

// GetLabelsInIssue gets labels attached on a single issue whose id is num.
func (c *Client) GetLabelsInIssue(num int) ([]*github.Label, error) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	labels, _, err := c.Client.Issues.ListLabelsByIssue(context.Background(), c.owner, c.repo, num, nil)
	if err != nil {
		logrus.Errorf("failed to get labels in issue %d: %v", num, err)
		return nil, err
	}
	logrus.Debugf("succeed in getting labels in issue %d", num)
	return labels, nil
}

// AddLabelsToIssue adds labels to an issue
func (c *Client) AddLabelsToIssue(num int, labels []string) error {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	if _, _, err := c.Client.Issues.AddLabelsToIssue(context.Background(), c.owner, c.repo, num, labels); err != nil {
		logrus.Errorf("failed to add labels %s to issue(pr) %d: %v", labels, num, err)
		return err
	}
	logrus.Debugf("succeed in adding labels %v for issue %d", labels, num)
	return nil
}

// RemoveLabelForIssue removes a label from an issue.
func (c *Client) RemoveLabelForIssue(num int, label string) error {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	if _, err := c.Client.Issues.RemoveLabelForIssue(context.Background(), c.owner, c.repo, num, label); err != nil {
		logrus.Errorf("failed to remove label %s for issue(pr) %d: %v", label, num, err)
		return err
	}
	logrus.Debugf("succeed in removing label %v for issue %d", label, num)
	return nil
}

// ReplaceLabelsForIssue replaces all labels for an issue.
func (c *Client) ReplaceLabelsForIssue(num int, labels []string) error {
	if _, _, err := c.Client.Issues.ReplaceLabelsForIssue(context.Background(), c.owner, c.repo, num, labels); err != nil {
		logrus.Errorf("failed to replace labels %v for issue(pr) %d: %v", labels, num, err)
		return err
	}
	logrus.Debugf("succeed in replacing labels %v for issue %d", labels, num)
	return nil
}

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
	logrus.Debugf("succeed in removing comment %s for pull request", id)
	return nil
}

// AssignIssueToUsers assigns users to the specified issue.
func (c *Client) AssignIssueToUsers(num int, users []string) error {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	if _, _, err := c.Client.Issues.AddAssignees(context.Background(), c.owner, c.repo, num, users); err != nil {
		logrus.Errorf("failed to assign users %s to issue(pr) %d: %v", users, num, err)
		return err
	}
	logrus.Debugf("succeed in assign users %s for pull request %d", users, num)
	return nil
}

// UnassignIssueToUsers assigns users to the specified issue.
func (c *Client) UnassignIssueToUsers(num int, users []string) error {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	if _, _, err := c.Client.Issues.AddAssignees(context.Background(), c.owner, c.repo, num, users); err != nil {
		logrus.Errorf("failed to assign users %s to issue(pr) %d: %v", users, num, err)
		return err
	}
	logrus.Debugf("succeed in assign users %s for pull request %d", users, num)
	return nil
}

// IssueHasLabel judges if an issue has a specified label.
func (c *Client) IssueHasLabel(num int, inputLabel string) bool {
	labels, err := c.GetLabelsInIssue(num)
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

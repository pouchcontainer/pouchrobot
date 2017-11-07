package issueCommentProcessor

import (
	"strings"

	"github.com/google/go-github/github"
)

// ActToIssueCommentCreated acts to issue comment.
// It covers the following parts:
// assign to user if he comments `#dibs`
func (icp *IssueCommentProcessor) ActToIssueCommentCreated(issue *github.Issue, comment *github.IssueComment) error {
	if comment.Body == nil || comment.User == nil || comment.User.Login == nil {
		return nil
	}

	commentUser := *(comment.User.Login)
	commentBody := *(comment.Body)

	users := []string{commentUser}

	if !strings.Contains(strings.ToLower(commentBody), "#dibs") {
		return nil
	}

	return icp.Client.AssignIssueToUsers(*(issue.Number), users)
}

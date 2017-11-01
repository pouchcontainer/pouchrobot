package issueCommentProcessor

import (
	"strings"

	"github.com/allencloud/automan/server/gh"
	"github.com/allencloud/automan/server/utils"
	"github.com/google/go-github/github"
)

// IssueCommentProcessor is
type IssueCommentProcessor struct {
	Client *gh.Client
}

// Process processes issue comment events
func (icp *IssueCommentProcessor) Process(data []byte) error {
	// process details
	actionType, err := utils.ExtractActionType(data)
	if err != nil {
		return err
	}

	issue, err := utils.ExactIssue(data)
	if err != nil {
		return err
	}

	comment, err := utils.ExactIssueComment(data)
	if err != nil {
		return err
	}

	switch actionType {
	case "created", "edited":
		if err := icp.ActToIssueCommentCreated(&issue, &comment); err != nil {
			return err
		}
	case "deleted":
		if err := icp.ActToIssueCommentDeleted(&issue, &comment); err != nil {
			return err
		}
	case "review_requested":
	}
	return nil
}

// ActToIssueCommentCreated acts to issue comment.
// It covers the following parts:
// assign to user if he comments `#dibs`
func (icp *IssueCommentProcessor) ActToIssueCommentCreated(issue *github.Issue, comment *github.IssueComment) error {
	if comment == nil || issue == nil {
		return nil
	}

	commentUser := *(comment.User.Login)
	commentBody := *(comment.Body)
	users := []string{commentUser}

	if strings.Contains(strings.ToLower(commentBody), "#dibs") {
		if err := icp.Client.AssignIssueToUsers(*(issue.Number), users); err != nil {
			return err
		}
	}

	return nil
}

// ActToIssueCommentDeleted acts an event that an issue comment is deleted.
func (icp *IssueCommentProcessor) ActToIssueCommentDeleted(issue *github.Issue, comment *github.IssueComment) error {
	if comment == nil || issue == nil {
		return nil
	}

	return nil
}

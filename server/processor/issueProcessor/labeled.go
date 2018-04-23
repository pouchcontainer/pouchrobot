package issueProcessor

import (
	"fmt"

	"github.com/google/go-github/github"
	"github.com/pouchcontainer/pouchrobot/server/utils"
)

// ActToIssueLabeled acts to issue labeled events
func (ip *IssueProcessor) ActToIssueLabeled(issue *github.Issue) error {
	ip.actToPriority(issue)
	return nil
}

func (ip *IssueProcessor) actToPriority(issue *github.Issue) error {
	if !ip.Client.IssueHasLabel(*(issue.Number), utils.PriorityP1Label) {
		id, exist := ip.Client.IssueHasComment(*(issue.Number), utils.IssueNeedP1CommentSubStr)
		if !exist {
			return nil
		}
		return ip.Client.RemoveComment(id)
	}

	if _, exist := ip.Client.IssueHasComment(*(issue.Number), utils.IssueNeedP1CommentSubStr); exist {
		return nil
	}

	body := fmt.Sprintf(utils.IssueNeedP1Comment, *(issue.User.Login))
	newComment := &github.IssueComment{
		Body: &body,
	}

	return ip.Client.AddCommentToIssue(*(issue.Number), newComment)
}

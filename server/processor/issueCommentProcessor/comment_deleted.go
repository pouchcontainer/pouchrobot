package issueCommentProcessor

import "github.com/google/go-github/github"

// ActToIssueCommentDeleted acts an event that an issue comment is deleted.
func (icp *IssueCommentProcessor) ActToIssueCommentDeleted(issue *github.Issue, comment *github.IssueComment) error {
	if comment == nil || issue == nil {
		return nil
	}

	return nil
}

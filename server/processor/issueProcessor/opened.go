package issueProcessor

import (
	"fmt"

	"github.com/google/go-github/github"
	"github.com/pouchcontainer/pouchrobot/server/processor/issueProcessor/open"
	"github.com/pouchcontainer/pouchrobot/server/utils"
)

// ActToIssueOpened acts to opened issue
// This function covers the following part:
// generate labels;
// attach comments;
// assign issue to specific user;
func (ip *IssueProcessor) ActToIssueOpened(issue *github.Issue) error {
	ip.attachLabels(issue)
	ip.attachComments(issue)

	return nil
}

func (ip *IssueProcessor) attachLabels(issue *github.Issue) error {
	labels := open.ParseToGenerateLabels(issue)
	if len(labels) == 0 {
		return nil
	}
	// TODO: check versions to add labels
	// only labels generated do we attach labels to issue
	return ip.Client.AddLabelsToIssue(*(issue.Number), labels)
}

func (ip *IssueProcessor) attachComments(issue *github.Issue) error {
	ip.attachTitleComments(issue)
	ip.attachBodyComments(issue)

	return nil
}

func (ip *IssueProcessor) attachTitleComments(issue *github.Issue) error {
	// check if the title is too short or the body empty.
	if issue.Title != nil && len(*(issue.Title)) > 20 {
		return nil
	}

	// attach comment
	body := fmt.Sprintf(utils.IssueTitleTooShort, *(issue.User.Login))
	newComment := &github.IssueComment{
		Body: &body,
	}

	if err := ip.Client.AddCommentToIssue(*(issue.Number), newComment); err != nil {
		return err
	}

	labels := []string{"status/more-info-needed"}
	return ip.Client.AddLabelsToIssue(*(issue.Number), labels)
}

func (ip *IssueProcessor) attachBodyComments(issue *github.Issue) error {
	if issue.Body != nil && len(*(issue.Body)) > 100 {
		return nil
	}

	// attach comment
	body := fmt.Sprintf(utils.IssueDescriptionTooShort, *(issue.User.Login))
	newComment := &github.IssueComment{
		Body: &body,
	}
	if err := ip.Client.AddCommentToIssue(*(issue.Number), newComment); err != nil {
		return err
	}

	if ip.Client.IssueHasLabel(*(issue.Number), "status/more-info-needed") {
		return nil
	}

	labels := []string{"status/more-info-needed"}
	return ip.Client.AddLabelsToIssue(*(issue.Number), labels)
}

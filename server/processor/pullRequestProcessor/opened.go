package pullRequestProcessor

import (
	"fmt"

	"github.com/allencloud/automan/server/processor/pullRequestProcessor/open"
	"github.com/allencloud/automan/server/utils"
	"github.com/google/go-github/github"
)

// ActToPROpened acts a pull request opened event.
func (prp *PullRequestProcessor) ActToPROpened(pr *github.PullRequest) error {
	prp.attachLabels(pr)
	prp.attachComments(pr)
	return nil
}

func (prp *PullRequestProcessor) attachLabels(pr *github.PullRequest) error {
	// attach labels
	labels := open.ParseToGeneratePRLabels(pr)
	if len(labels) == 0 {
		return nil
	}
	return prp.Client.AddLabelsToIssue(*(pr.Number), labels)
}

func (prp *PullRequestProcessor) attachComments(pr *github.PullRequest) error {
	prp.attachTitleComments(pr)
	prp.attachBodyComments(pr)
	if err := prp.addSignoffComments(pr); err != nil {
		return err
	}
	prp.attachFirstContributionComments(pr)
	return nil
}

func (prp *PullRequestProcessor) attachTitleComments(pr *github.PullRequest) error {
	if pr.Title != nil && len(*(pr.Title)) > 20 {
		return nil
	}

	// attach comment
	body := fmt.Sprintf(utils.PRTitleTooShort, *(pr.User.Login))
	newComment := &github.IssueComment{
		Body: &body,
	}

	return prp.Client.AddCommentToPR(*(pr.Number), newComment)
}

func (prp *PullRequestProcessor) attachBodyComments(pr *github.PullRequest) error {
	if pr.Body != nil && len(*(pr.Body)) > 50 {
		return nil
	}

	body := fmt.Sprintf(utils.PRDescriptionTooShort, *(pr.User.Login))
	newComment := &github.IssueComment{
		Body: &body,
	}

	return prp.Client.AddCommentToPR(*(pr.Number), newComment)
}

func (prp *PullRequestProcessor) addSignoffComments(pr *github.PullRequest) error {
	// check whether commits are following the rules
	commits, err := prp.Client.ListCommits(*(pr.Number))
	if err != nil {
		return err
	}

	needSignoff := false
	for _, commit := range commits {
		if commit.Commit != nil && !dcoRegex.MatchString(*commit.Commit.Message) {
			needSignoff = true
			break
		}
	}

	if !needSignoff {
		return nil
	}

	body := fmt.Sprintf(utils.PRNeedsSignOff, *(pr.User.Login))
	newComment := &github.IssueComment{
		Body: &body,
	}

	return prp.Client.AddCommentToPR(*(pr.Number), newComment)
}

func (prp *PullRequestProcessor) attachFirstContributionComments(pr *github.PullRequest) error {
	// check whether this is the first contributor of the committer
	if pr.AuthorAssociation == nil {
		return nil
	}

	if !isFirstContribution(*(pr.AuthorAssociation)) {
		return nil
	}

	body := fmt.Sprintf(utils.FirstCommitComment, *(pr.User.Login))
	newComment := &github.IssueComment{
		Body: &body,
	}
	return prp.Client.AddCommentToPR(*(pr.Number), newComment)
}

func isFirstContribution(str string) bool {
	// TODO: add how to check
	return false
}

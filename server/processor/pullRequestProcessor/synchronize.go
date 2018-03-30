package pullRequestProcessor

import (
	"context"
	"fmt"
	"strings"

	"github.com/pouchcontainer/pouchrobot/server/processor/pullRequestProcessor/open"
	"github.com/pouchcontainer/pouchrobot/server/utils"

	"github.com/google/go-github/github"
)

// ActToPRSynchronized acts to event that a pr is synchronized.
func (prp *PullRequestProcessor) ActToPRSynchronized(syncPR *github.PullRequest) error {
	prp.removeConflictLabel(syncPR)
	prp.changeSizeLabel(syncPR)
	prp.changeSignCommitComment(syncPR)
	return nil
}

func (prp *PullRequestProcessor) removeConflictLabel(syncPR *github.PullRequest) error {
	pr, err := prp.Client.GetSinglePR(*(syncPR.Number))
	if err != nil {
		return nil
	}

	// check if this pr is updated to solve the conflict,
	// if that remove label 'conflict/needs-rebase' and remove the relating comment.
	if !prp.Client.IssueHasLabel(*(pr.Number), utils.PRConflictLabel) {
		// pull request has no conflict label, do nothing
		return nil
	}

	if pr.Mergeable == nil || *(pr.Mergeable) == false {
		return nil
	}

	// remove conflict label
	prp.Client.RemoveLabelForIssue(*(pr.Number), utils.PRConflictLabel)
	// remove conflict comment
	prp.RemoveConflictComment(context.Background(), *(pr.Number))

	return nil
}

func (prp *PullRequestProcessor) changeSizeLabel(pr *github.PullRequest) error {
	// check if we need to change the PR size label
	newSizeLabel := open.ParseToGetPRSize(pr)

	if prp.Client.IssueHasLabel(*(pr.Number), newSizeLabel) {
		// pull request already has newSize label, do nothing
		return nil
	}

	// pull request has no newSizeLabel, do following things:
	// remove original size label
	// add newSizeLabel

	originalLabels, err := prp.Client.GetLabelsInIssue(*(pr.Number))
	if err != nil {
		return err
	}

	for _, label := range originalLabels {
		if strings.HasPrefix(*(label.Name), utils.SizeLabelPrefix) {
			prp.Client.RemoveLabelForIssue(*(pr.Number), label.GetName())
			break
		}
	}

	newLabels := []string{newSizeLabel}
	prp.Client.AddLabelsToIssue(*(pr.Number), newLabels)

	return nil
}

// RemoveConflictComment removes a conflict comment for a pull request
func (prp *PullRequestProcessor) RemoveConflictComment(ctx context.Context, num int) error {
	prComments, err := prp.Client.ListComments(num)
	if err != nil {
		return err
	}
	for _, comment := range prComments {
		commentBody := *(comment.Body)
		subBody := utils.PRConflictSubStr
		if strings.HasSuffix(commentBody, subBody) {
			// remove all if there are more than one
			prp.Client.RemoveComment(*(comment.ID))
		}
	}
	return nil
}

// changeSignCommitComment changes comments of being signed off.
func (prp *PullRequestProcessor) changeSignCommitComment(pr *github.PullRequest) error {
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

	// try to remove sign off commits if there are any.
	prp.Client.RmCommentsViaStr(*(pr.Number), utils.PRNeedsSignOffStr)

	if !needSignoff {
		return nil
	}

	body := fmt.Sprintf(utils.PRNeedsSignOff, *(pr.User.Login))
	newComment := &github.IssueComment{
		Body: &body,
	}

	return prp.Client.AddCommentToPR(*(pr.Number), newComment)
}

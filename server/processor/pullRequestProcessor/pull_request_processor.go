package pullRequestProcessor

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/allencloud/automan/server/gh"
	"github.com/allencloud/automan/server/processor/pullRequestProcessor/open"
	"github.com/allencloud/automan/server/utils"

	"github.com/google/go-github/github"
	"github.com/sirupsen/logrus"
)

var (
	// SizePrefix means the prefix of a size label
	SizePrefix = "SIZE/"
	dcoRegex   = regexp.MustCompile("(?m)(Docker-DCO-1.1-)?Signed-off-by: ([^<]+) <([^<>@]+@[^<>]+)>( \\(github: ([a-zA-Z0-9][a-zA-Z0-9-]+)\\))?")
)

// PullRequestProcessor is
type PullRequestProcessor struct {
	Client *gh.Client
}

// Process processes pull request events
func (prp *PullRequestProcessor) Process(data []byte) error {
	// process details
	actionType, err := utils.ExtractActionType(data)
	if err != nil {
		return err
	}

	logrus.Infof("received event type [pull request], action type [%s]", actionType)

	pr, err := utils.ExactPR(data)
	if err != nil {
		return err
	}
	logrus.Debugf("pull request: %v", pr)

	switch actionType {
	case "opened":
		if err := prp.ActToPROpenOrEdit(&pr); err != nil {
			return err
		}
	case "labeled":
		if err := prp.ActToPRLabeled(&pr); err != nil {
			return err
		}
	case "review_requested":
	case "synchronize":
		if err := prp.ActToPRSynchronized(&pr); err != nil {
			return err
		}
	case "edited":
		if err := prp.ActToPROpenOrEdit(&pr); err != nil {
			return err
		}
	case "pull_request_review":
	default:
		return fmt.Errorf("unknown action type %s in pull request: ", actionType)
	}
	return nil
}

// ActToPROpenOrEdit acts
func (prp *PullRequestProcessor) ActToPROpenOrEdit(pr *github.PullRequest) error {
	// attach labels
	labels := open.ParseToGeneratePRLabels(pr)
	if len(labels) != 0 {
		// only labels generated do we attach labels to issue
		if err := prp.Client.AddLabelsToIssue(*(pr.Number), labels); err != nil {
			return err
		}
	}

	// attach comment
	newComment := &github.IssueComment{}
	// check if the title is too short or the body empty.
	if pr.Title == nil || len(*(pr.Title)) < 20 {
		body := fmt.Sprintf(utils.PRTitleTooShort, *(pr.User.Login))
		newComment.Body = &body
		if err := prp.Client.AddCommentToPR(*(pr.Number), newComment); err != nil {
			return err
		}
		return nil
	}

	if pr.Body == nil || len(*(pr.Body)) < 50 {
		body := fmt.Sprintf(utils.PRDescriptionTooShort, *(pr.User.Login))
		newComment.Body = &body
		if err := prp.Client.AddCommentToPR(*(pr.Number), newComment); err != nil {
			return err
		}
		return nil
	}

	commits, err := prp.Client.ListCommits(*(pr.Number))
	if err != nil {
		return err
	}

	for _, commit := range commits {
		if commit.Commit != nil && !dcoRegex.MatchString(*commit.Commit.Message) {
			// pull request is not signed
			// TODO add implementation
			break
		}
	}

	return nil
}

// ActToPRLabeled acts the event of pull request labeled.
func (prp *PullRequestProcessor) ActToPRLabeled(pr *github.PullRequest) error {
	return nil
}

// ActToPRSynchronized acts to event that a pr is synchronized.
func (prp *PullRequestProcessor) ActToPRSynchronized(syncPR *github.PullRequest) error {
	pr, err := prp.Client.GetSinglePR(*(syncPR.Number))
	if err != nil {
		return nil
	}
	// check if this pr is updated to solve the conflict,
	// if that remove label 'conflict/needs-rebase' and remove the relating comment.
	if prp.Client.IssueHasLabel(*(pr.Number), utils.ConflictLabel) {
		if pr.Mergeable != nil && *(pr.Mergeable) == true {
			// remove conflict label
			prp.Client.RemoveLabelForIssue(*(pr.Number), utils.ConflictLabel)

			// remove conflict comment
			prp.RemoveConflictComment(context.Background(), *(pr.Number))
		}
	}

	// check if we need to change the PR size label
	newSizeLabel := open.ParseToGetPRSize(pr)
	if !prp.Client.IssueHasLabel(*(pr.Number), newSizeLabel) {
		// first remove the original size label
		originalLabels, err := prp.Client.GetLabelsInIssue(*(pr.Number))
		if err != nil {
			return err
		}
		for _, label := range originalLabels {
			if strings.HasPrefix(label.GetName(), SizePrefix) {
				prp.Client.RemoveLabelForIssue(*(pr.Number), label.GetName())
				break
			}
		}
		newLabels := []string{newSizeLabel}
		prp.Client.AddLabelsToIssue(*(pr.Number), newLabels)
	}

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
		subBody := utils.ConflictSubStr
		if strings.HasSuffix(commentBody, subBody) {
			return prp.Client.RemoveComment(*(comment.ID))
		}
	}
	return nil
}

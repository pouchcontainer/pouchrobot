package pullRequestProcessor

import (
	"context"
	"fmt"
	"strings"

	"github.com/allencloud/automan/server/gh"
	"github.com/allencloud/automan/server/processor/pullRequestProcessor/open"
	putils "github.com/allencloud/automan/server/processor/utils"
	"github.com/allencloud/automan/server/utils"
	"github.com/google/go-github/github"
	"github.com/sirupsen/logrus"
)

var (
	// ConflictLabel means a label
	ConflictLabel = "conflict/needs-rebase"
	// SizePrefix means the prefix of a size label
	SizePrefix = "SIZE/"
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

	switch actionType {
	case "opened":
		if err := prp.ActToPROpenOrEdit(&pr); err != nil {
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
		if err := prp.Client.AddLabelsToIssue(context.Background(), *(pr.Number), labels); err != nil {
			return err
		}
	}

	// attach comment
	newComment := &github.IssueComment{}
	// check if the title is too short or the body empty.
	if pr.Title == nil || len(*(pr.Title)) < 20 {
		body := fmt.Sprintf(putils.PRTitleTooShort, *(pr.User.Login))
		newComment.Body = &body
		if err := prp.Client.AddCommentToPR(context.Background(), *(pr.Number), newComment); err != nil {
			return err
		}
		return nil
	}

	if pr.Body == nil || len(*(pr.Body)) < 50 {
		body := fmt.Sprintf(putils.PRDescriptionTooShort, *(pr.User.Login))
		newComment.Body = &body
		if err := prp.Client.AddCommentToPR(context.Background(), *(pr.Number), newComment); err != nil {
			return err
		}
		return nil
	}
	return nil
}

// ActToPRSynchronized acts to event that a pr is synchronized.
func (prp *PullRequestProcessor) ActToPRSynchronized(pr *github.PullRequest) error {
	// check if this pr is updated to solve the conflict,
	// if that remove label 'conflict/needs-rebase' and remove the relating comment.
	if prp.Client.IssueHasLabel(*(pr.Number), ConflictLabel) {
		if pr.Mergeable != nil && *(pr.Mergeable) == true {
			// remove conflict label
			prp.Client.RemoveLabelForIssue(context.Background(), *(pr.Number), ConflictLabel)

			// remove conflict comment
			prp.RemoveConflictComment(context.Background(), *(pr.Number))
		}
	}

	// check if we need to change the PR size label
	newSizeLabel := open.ParseToGetPRSize(pr)
	if !prp.Client.IssueHasLabel(*(pr.Number), newSizeLabel) {
		// first remove the original size label
		originalLabels, err := prp.Client.GetLabelsInIssue(context.Background(), *(pr.Number))
		if err != nil {
			return err
		}
		for _, label := range originalLabels {
			if strings.HasPrefix(label.GetName(), SizePrefix) {
				prp.Client.RemoveLabelForIssue(context.Background(), *(pr.Number), label.GetName())
				break
			}
		}
		newLabels := []string{newSizeLabel}
		prp.Client.AddLabelsToIssue(context.Background(), *(pr.Number), newLabels)
	}

	return nil
}

// RemoveConflictComment removes a conflict comment for a pull request
func (prp *PullRequestProcessor) RemoveConflictComment(ctx context.Context, num int) error {
	prComments, err := prp.Client.ListPRComments(context.Background(), num)
	if err != nil {
		return err
	}
	for _, comment := range prComments {
		commentBody := comment.GetBody()
		subBody := `
Conflict happens after merging a previous commit.
Please rebase the branch against master and push it again.
Thanks a lot.`
		if strings.HasSuffix(commentBody, subBody) {
			return prp.Client.RemoveCommentForPR(context.Background(), *(comment.ID))
		}
	}
	return nil
}

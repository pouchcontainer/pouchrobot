package pullRequestProcessor

import (
	"context"
	"fmt"

	"github.com/allencloud/automan/server/gh"
	"github.com/allencloud/automan/server/processor/pullRequestProcessor/open"
	putils "github.com/allencloud/automan/server/processor/utils"
	"github.com/allencloud/automan/server/utils"
	"github.com/google/go-github/github"
	"github.com/sirupsen/logrus"
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

	logrus.Infof("Received a PR: %v", pr)

	switch actionType {
	case "opened":
		if err := prp.ActToPROpenOrEdit(&pr); err != nil {
			return err
		}
	case "review_requested":
	case "synchronize":
		logrus.Info("-------------------------------------Got a synchronized event")
		logrus.Infof("the mergeable is %s", *(pr.Mergeable))
	case "edited":
		if err := prp.ActToPROpenOrEdit(&pr); err != nil {
			return err
		}
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
		if err := prp.Client.AddLabelsToPR(context.Background(), *(pr.Number), labels); err != nil {
			logrus.Errorf("failed to add labels %v to issue %d: %v", labels, *(pr.Number), err)
			return err
		}
		logrus.Infof("succeed in attaching labels %v to issue %d", labels, *(pr.Number))
	}

	// attach comment
	newComment := &github.PullRequestComment{}

	// check if the title is too short or the body empty.
	if len(*(pr.Title)) < 20 {
		body := fmt.Sprintf(putils.PRTitleTooShort, *(pr.User.Login))
		newComment.Body = &body
		if err := prp.Client.AddCommentToPR(context.Background(), *(pr.Number), newComment); err != nil {
			logrus.Errorf("failed to add TITLE TOO SHORT comment to pr %d", *(pr.Number))
			return err
		}
		logrus.Infof("succeed in attaching TITLE TOO SHORT comment for pr %d", *(pr.Number))

		return nil
	}

	if pr.Body == nil || *(pr.Body) == "" || len(*(pr.Body)) < 50 {
		body := fmt.Sprintf(putils.PRDescriptionTooShort, *(pr.User.Login))
		newComment.Body = &body
		if err := prp.Client.AddCommentToPR(context.Background(), *(pr.Number), newComment); err != nil {
			logrus.Errorf("failed to add EMPTY OR TOO SHORT PR BODY comment to pr %d", *(pr.Number))
			return err
		}
		logrus.Infof("succeed in attaching BODY TOO SHORT comment for pr %d", *(pr.Number))
		return nil
	}
	return nil
}

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

	issue, err := utils.ExactIssue(data)
	if err != nil {
		return err
	}

	comment, err := utils.ExactIssueComment(data)
	if err != nil {
		return nil
	}

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
		logrus.Info("-------------------------------------Got a synchronized event")
		logrus.Infof("the mergeable is %v", pr.Mergeable)
	case "edited":
		if err := prp.ActToPROpenOrEdit(&pr); err != nil {
			return err
		}
	case "pull_request_review":
	case "created":
		logrus.Infof("Got an issue: %v", issue)
		if err := prp.ActToPRCommented(&issue, &comment); err != nil {
			return nil
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

// ActToPRCommented acts added comment to the PR
// Here are the rules:
// 1. if maintainers attached LGTM and currently no LGTM, add a label "LGTM";
// 2. if maintainers attached LGTM and already has a LGTM, add a label "APPROVED";
func (prp *PullRequestProcessor) ActToPRCommented(issue *github.Issue, comment *github.IssueComment) error {
	body := *(comment.Body)
	user := *(issue.User.Login)
	logrus.Infof("body: %s, user:%s", body, user)
	if hasLGTMFromMaintainer(user, body) && noLGTMInLabels(issue) {
		prp.Client.AddLabelsToPR(context.Background(), *(issue.Number), []string{"LGTM"})
	}
	if hasLGTMFromMaintainer(user, body) && hasLGTMInLabels(issue) {
		prp.Client.AddLabelsToPR(context.Background(), *(issue.Number), []string{"APPROVED"})
	}
	return nil
}

func hasLGTMFromMaintainer(user string, body string) bool {
	if !strings.Contains(strings.ToLower(body), "lgtm") {
		return false
	}

	for _, maintainerID := range putils.Maintainers {
		if strings.ToLower(user) == strings.ToLower(maintainerID) {
			return true
		}
	}
	return false
}

func noLGTMInLabels(issue *github.Issue) bool {
	for _, label := range issue.Labels {
		if label.GetName() == "LGTM" {
			return false
		}
	}
	return true
}

func hasLGTMInLabels(issue *github.Issue) bool {
	for _, label := range issue.Labels {
		if label.GetName() == "LGTM" {
			return true
		}
	}
	return false
}

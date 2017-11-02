package prCommentProcessor

import (
	"strings"

	"github.com/allencloud/automan/server/gh"
	putils "github.com/allencloud/automan/server/processor/utils"
	"github.com/allencloud/automan/server/utils"
	"github.com/google/go-github/github"
	"github.com/sirupsen/logrus"
)

// PRCommentProcessor is
type PRCommentProcessor struct {
	Client *gh.Client
}

// Process processes pull request events
func (prcp *PRCommentProcessor) Process(data []byte) error {
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
	case "created":
		//logrus.Infof("Got an issue: %v", issue)
		if err := prcp.ActToPRCommented(&issue, &comment); err != nil {
			return nil
		}
	case "edited":
	}
	return nil
}

// ActToPRCommented acts added comment to the PR
// Here are the rules:
// 1. if maintainers attached LGTM and currently no LGTM, add a label "LGTM";
// 2. if maintainers attached LGTM and already has a LGTM, add a label "APPROVED";
func (prcp *PRCommentProcessor) ActToPRCommented(issue *github.Issue, comment *github.IssueComment) error {
	body := *(comment.Body)
	user := *(issue.User.Login)
	logrus.Infof("body: %s, user:%s, issue: %v", body, user, *issue)

	if hasLGTMFromMaintainer(user, body) && noLGTMInLabels(issue) {
		prcp.Client.AddLabelsToIssue(*(issue.Number), []string{"LGTM"})
	}
	// FIXME: one maintainer attached two LGTMs, it will attach an APPROVED
	if hasLGTMFromMaintainer(user, body) && hasLGTMInLabels(issue) {
		prcp.Client.AddLabelsToIssue(*(issue.Number), []string{"APPROVED"})
	}
	return nil
}

func hasLGTMFromMaintainer(user string, body string) bool {
	// FIXME: if a maintainer attached a comment like: LGTM if change request solved,
	// this rule will still add a lgtm label
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

package issueProcessor

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/allencloud/automan/server/gh"
	"github.com/allencloud/automan/server/processor/issueProcessor/open"
	"github.com/google/go-github/github"
)

// IssueOpenedActionFunc defines
type IssueOpenedActionFunc func(issue *github.Issue) []string

// IssueOpenedActionFuncs defines
var IssueOpenedActionFuncs []IssueOpenedActionFunc

// Register registers IssueOpenedActionFunc
func Register(a IssueOpenedActionFunc) {
	IssueOpenedActionFuncs = append(IssueOpenedActionFuncs, a)
}

func init() {
	funcs := []IssueOpenedActionFunc{
		open.ParseToGenerateLabels,
	}

	for _, processFunc := range funcs {
		Register(processFunc)
	}
}

// TriggeredIssueProcessor is
type TriggeredIssueProcessor struct {
	Client *gh.Client
}

// Process processes
func (fIP *TriggeredIssueProcessor) Process(data []byte) error {
	// process details
	actionType, err := ExtractActionType(data)
	if err != nil {
		return err
	}

	logrus.Infof("received event type [issues], action type [%s]", actionType)

	issue, err := ExactIssue(data)
	if err != nil {
		return err
	}

	switch actionType {
	case "opened":
		if err := fIP.ActToIssueOpen(&issue); err != nil {
			return err
		}
	case "reopened":
	case "edited":
	default:
		return fmt.Errorf("unknown action type %s in issue: ", actionType)
	}

	return nil
}

// ActToIssueOpen acts to opened issue
func (fIP *TriggeredIssueProcessor) ActToIssueOpen(issue *github.Issue) error {
	// generate labels
	labels := open.ParseToGenerateLabels(issue)
	if err := fIP.Client.AddLabelsToIssue(context.Background(), *(issue.Number), labels); err != nil {
		logrus.Errorf("failed to add labels %v to issue %d: %v", labels, *(issue.Number), err)
		return err
	}
	logrus.Infof("succeed in attaching labels %v to issue %d", labels, *(issue.Number))

	return nil
}

// ExtractActionType extracts the action type.
func ExtractActionType(data []byte) (string, error) {
	var m struct {
		Action string `json:"action"`
	}
	if err := json.Unmarshal(data, &m); err != nil {
		return "", err
	}
	return m.Action, nil
}

// ExactIssue extracts the issue.
func ExactIssue(data []byte) (github.Issue, error) {
	var m struct {
		Issue github.Issue `json:"issue"`
	}
	if err := json.Unmarshal(data, &m); err != nil {
		return github.Issue{}, err
	}
	return m.Issue, nil
}

func exactPR(data []byte) (github.PullRequest, error) {
	var m struct {
		PullRequest github.PullRequest `json:"pull_reques"`
	}
	if err := json.Unmarshal(data, &m); err != nil {
		return github.PullRequest{}, err
	}
	return m.PullRequest, nil
}

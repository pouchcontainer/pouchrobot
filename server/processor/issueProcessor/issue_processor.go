package issueProcessor

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/google/go-github/github"
	"github.com/pouchcontainer/pouchrobot/server/gh"
	"github.com/pouchcontainer/pouchrobot/server/processor/issueProcessor/open"
	"github.com/pouchcontainer/pouchrobot/server/utils"
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

// IssueProcessor is
type IssueProcessor struct {
	Client *gh.Client
}

// Process processes
func (ip *IssueProcessor) Process(data []byte) error {
	// process details
	actionType, err := utils.ExtractActionType(data)
	if err != nil {
		return err
	}

	logrus.Infof("received event type [issues], action type [%s]", actionType)

	issue, err := utils.ExactIssue(data)
	if err != nil {
		return err
	}

	switch actionType {
	case "opened":
		if err := ip.ActToIssueOpened(&issue); err != nil {
			return err
		}
	case "edited":
		if err := ip.ActToIssueEdited(&issue); err != nil {
			return err
		}
	case "labeled":
		if err := ip.ActToIssueLabeled(&issue); err != nil {
			return nil
		}
	case "reopened":
	default:
		return fmt.Errorf("unknown action type %s in issue: ", actionType)
	}

	return nil
}

package issueCommentProcessor

import (
	"github.com/pouchcontainer/pouchrobot/server/gh"
	"github.com/pouchcontainer/pouchrobot/server/utils"
)

// IssueCommentProcessor is
type IssueCommentProcessor struct {
	Client *gh.Client
}

// Process processes issue comment events
func (icp *IssueCommentProcessor) Process(data []byte) error {
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
	case "created", "edited":
		if err := icp.ActToIssueCommentCreated(&issue, &comment); err != nil {
			return err
		}
	case "deleted":
		if err := icp.ActToIssueCommentDeleted(&issue, &comment); err != nil {
			return err
		}
	case "review_requested":
	}
	return nil
}

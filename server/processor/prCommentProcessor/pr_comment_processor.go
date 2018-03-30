package prCommentProcessor

import (
	"github.com/pouchcontainer/pouchrobot/server/gh"
	"github.com/pouchcontainer/pouchrobot/server/utils"
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
		if err := prcp.ActToPRCommented(&issue, &comment); err != nil {
			return nil
		}
	case "edited":
		if err := prcp.ActToPRCommentEdited(&issue, &comment); err != nil {
			return nil
		}
	}
	return nil
}

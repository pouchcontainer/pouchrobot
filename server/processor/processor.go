package processor

import (
	"fmt"

	"github.com/allencloud/automan/server/gh"

	"github.com/allencloud/automan/server/processor/issueProcessor"
	"github.com/allencloud/automan/server/processor/pullRequestProcessor"
	"github.com/allencloud/automan/server/utils"
)

type processor interface {
	// Process processes item automan gets, and then execute operations torwards items on GitHub
	Process(data []byte) error
}

// Processor contains several specific processors
type Processor struct {
	IssueProcessor       *issueProcessor.TriggeredIssueProcessor
	PullRequestProcessor *pullRequestProcessor.PullRequestProcessor
}

// NewProcessor creates
func NewProcessor(client *gh.Client) *Processor {
	return &Processor{
		IssueProcessor: &issueProcessor.TriggeredIssueProcessor{
			Client: client,
		},
		PullRequestProcessor: &pullRequestProcessor.PullRequestProcessor{
			Client: client,
		},
	}
}

// HandleEvent processes an event received from github
func (p *Processor) HandleEvent(eventType string, data []byte) error {
	// since pr is also a kind of issue, we need to first make it clear
	issueType := judgeIssueOrPR(data)
	if issueType == "issue" {
		eventType = "issues"
	} else if issueType == "pull_request" {
		eventType = "pull_request"
	}

	switch eventType {
	case "issues":
		p.IssueProcessor.Process(data)
	case "issue_comment":
		p.IssueProcessor.Process(data)
	case "pull_request":
		p.PullRequestProcessor.Process(data)
	default:
		return fmt.Errorf("unknown event type %s", eventType)
	}
	return nil
}

func judgeIssueOrPR(data []byte) string {
	issue, err := utils.ExactIssue(data)
	if err != nil {
		return ""
	}

	if issue.PullRequestLinks == nil {
		return "issue"
	}
	return "pull_request"
}

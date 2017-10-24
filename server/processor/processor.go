package processor

import (
	"fmt"

	"github.com/allencloud/automan/server/gh"

	"github.com/allencloud/automan/server/processor/issueProcessor"
	"github.com/allencloud/automan/server/processor/pullRequestProcessor"
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

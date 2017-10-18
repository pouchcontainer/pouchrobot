package processor

import (
	"fmt"

	"github.com/allencloud/automan/server/gh"

	"github.com/allencloud/automan/server/processor/issueProcessor"
	"github.com/sirupsen/logrus"
)

type processor interface {
	// Process processes item automan gets, and then execute operations torwards items on GitHub
	Process(data []byte) error
}

// Processor contains several specific processors
type Processor struct {
	IssueProcessor *issueProcessor.TriggeredIssueProcessor
}

// NewProcessor creates
func NewProcessor(client *gh.Client) *Processor {
	return &Processor{
		IssueProcessor: &issueProcessor.TriggeredIssueProcessor{
			Client: client,
		},
	}
}

// HandleEvent processes an event received from github
func (p *Processor) HandleEvent(eventType string, data []byte) error {
	logrus.Info("")
	switch eventType {
	case "issues":
		p.IssueProcessor.Process(data)
	case "pull_request":
		//processPullRequestEvent(data)
	default:
		return fmt.Errorf("unknown event type %s", eventType)
	}
	return nil
}

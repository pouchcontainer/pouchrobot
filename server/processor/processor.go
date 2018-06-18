// Copyright 2018 The Pouch Robot Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package processor

import (
	"fmt"

	"github.com/pouchcontainer/pouchrobot/server/processor/issueCommentProcessor"
	"github.com/pouchcontainer/pouchrobot/server/processor/issueProcessor"
	"github.com/pouchcontainer/pouchrobot/server/processor/prCommentProcessor"
	"github.com/pouchcontainer/pouchrobot/server/processor/pullRequestProcessor"
	"github.com/pouchcontainer/pouchrobot/server/utils"

	"github.com/pouchcontainer/pouchrobot/server/gh"
	"github.com/sirupsen/logrus"
)

type processor interface {
	// Process processes item automan gets, and then execute operations torwards items on GitHub
	Process(data []byte) error
}

// Processor contains several specific processors
type Processor struct {
	IssueProcessor        *issueProcessor.IssueProcessor
	PullRequestProcessor  *pullRequestProcessor.PullRequestProcessor
	IssueCommentProcessor *issueCommentProcessor.IssueCommentProcessor
	PRCommentProcessor    *prCommentProcessor.PRCommentProcessor
}

// New initializes a brand new processor.
func New(client *gh.Client) *Processor {
	return &Processor{
		IssueProcessor: &issueProcessor.IssueProcessor{
			Client: client,
		},
		PullRequestProcessor: &pullRequestProcessor.PullRequestProcessor{
			Client: client,
		},
		IssueCommentProcessor: &issueCommentProcessor.IssueCommentProcessor{
			Client: client,
		},
		PRCommentProcessor: &prCommentProcessor.PRCommentProcessor{
			Client: client,
		},
	}
}

// HandleEvent processes an event received from github
func (p *Processor) HandleEvent(eventType string, data []byte) error {
	logrus.Infof("eventType is %v", eventType)
	switch eventType {
	case "issues":
		p.IssueProcessor.Process(data)
	case "pull_request":
		p.PullRequestProcessor.Process(data)
	case "issue_comment":
		// since pr is also a kind of issue, we need to first make it clear
		issueType := judgeIssueOrPR(data)
		logrus.Infof("get issueType: %s", issueType)
		if issueType == "issue" {
			p.IssueCommentProcessor.Process(data)
			return nil
		}
		if issueType == "pull_request" {
			p.PRCommentProcessor.Process(data)
			return nil
		}
	case "ping":
		logrus.Debug("Got ping from GitHub")
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

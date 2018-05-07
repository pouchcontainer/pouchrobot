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

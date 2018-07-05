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

package issueCommentProcessor

import (
	"github.com/pouchcontainer/pouchrobot/gh"
	"github.com/pouchcontainer/pouchrobot/utils"
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

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

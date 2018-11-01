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
	"strings"

	"github.com/google/go-github/github"
)

// ActToIssueCommentCreated acts to issue comment.
// It covers the following parts:
// assign to user if he comments `#dibs` or `/assign`
func (icp *IssueCommentProcessor) ActToIssueCommentCreated(issue *github.Issue, comment *github.IssueComment) error {
	if comment.Body == nil || comment.User == nil || comment.User.Login == nil {
		return nil
	}

	commentUser := *(comment.User.Login)
	commentBody := *(comment.Body)

	users := []string{commentUser}

	if strings.HasPrefix(strings.ToLower(commentBody), "#dibs") || strings.HasPrefix(strings.ToLower(commentBody), "/assign") {
		return icp.Client.AssignIssueToUsers(*(issue.Number), users)
	}

	return nil
}

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

	"github.com/google/go-github/github"
	"github.com/pouchcontainer/pouchrobot/server/utils"
)

// ActToIssueLabeled acts to issue labeled events
func (ip *IssueProcessor) ActToIssueLabeled(issue *github.Issue) error {
	ip.actToPriority(issue)
	return nil
}

func (ip *IssueProcessor) actToPriority(issue *github.Issue) error {
	if !ip.Client.IssueHasLabel(*(issue.Number), utils.PriorityP1Label) {
		id, exist := ip.Client.IssueHasComment(*(issue.Number), utils.IssueNeedP1CommentSubStr)
		if !exist {
			return nil
		}
		return ip.Client.RemoveComment(id)
	}

	if _, exist := ip.Client.IssueHasComment(*(issue.Number), utils.IssueNeedP1CommentSubStr); exist {
		return nil
	}

	body := fmt.Sprintf(utils.IssueNeedP1Comment, *(issue.User.Login))
	newComment := &github.IssueComment{
		Body: &body,
	}

	return ip.Client.AddCommentToIssue(*(issue.Number), newComment)
}

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
	"strings"

	"github.com/google/go-github/github"
	"github.com/pouchcontainer/pouchrobot/server/utils"
)

// ActToPRCommented acts added comment to the PR
// Here are the rules:
func (prcp *PRCommentProcessor) ActToPRCommented(issue *github.Issue, comment *github.IssueComment) error {
	prcp.updateLabels(issue, comment)
	// retrigger test case by adding a comment of "/retest"
	prcp.retriggerCI(issue, comment)

	return nil
}

func (prcp *PRCommentProcessor) updateLabels(issue *github.Issue, comment *github.IssueComment) error {
	if comment.Body == nil || comment.User == nil || comment.User.Login == nil {
		return nil
	}

	body := *(comment.Body)
	user := *(comment.User.Login)

	if !strings.HasPrefix(strings.ToLower(body), "lgtm") && !strings.HasSuffix(strings.ToLower(body), "lgtm") {
		return nil
	}

	if !isMaintainer(user) {
		return nil
	}

	if prcp.Client.IssueHasLabel(*(issue.Number), "LGTM") {
		return nil
	}

	return prcp.Client.AddLabelsToIssue(*(issue.Number), []string{"LGTM"})
}

func (prcp *PRCommentProcessor) retriggerCI(issue *github.Issue, comment *github.IssueComment) error {
	if comment.Body == nil || comment.User == nil || comment.User.Login == nil {
		return nil
	}

	body := *(comment.Body)

	if !strings.Contains(body, "/retest") {
		return nil
	}
	// TODO: call CI system to retest
	return nil
}

func isMaintainer(user string) bool {
	for _, maintainerID := range utils.Maintainers {
		if strings.ToLower(user) == strings.ToLower(maintainerID) {
			return true
		}
	}
	return false
}

func hasLGTMInLabels(issue *github.Issue) bool {
	for _, label := range issue.Labels {
		if label.GetName() == "LGTM" {
			return true
		}
	}
	return false
}

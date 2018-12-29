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

	"github.com/pouchcontainer/pouchrobot/processor/issueProcessor/open"
	"github.com/pouchcontainer/pouchrobot/utils"

	"github.com/google/go-github/github"
)

// ActToIssueOpened acts to opened issue
// This function covers the following part:
// generate labels;
// attach comments;
// assign issue to specific user;
func (ip *IssueProcessor) ActToIssueOpened(issue *github.Issue) error {
	ip.attachLabels(issue)
	ip.attachComments(issue)
	ip.autoTranslate(issue)
	return nil
}

func (ip *IssueProcessor) autoTranslate(issue *github.Issue) error {
	translateTitle := ip.Translator.Translate(*issue.Title, false)
	translateBody := ip.Translator.Translate(*issue.Body, true)
	if translateTitle == "" && translateBody == "" {
		return nil
	}
	newIssue := &github.IssueRequest{}
	if translateTitle != "" {
		newIssue.Title = &translateTitle
	}
	if translateBody != "" {
		translateBody += "\r\n\r\n***!!!!WE STRONGLY ENCOURAGE YOU TO DESCRIBE YOUR ISSUE IN ENGLISH!!!!***"
		newIssue.Body = &translateBody
	}
	return ip.Client.EditIssue(*issue.Number, newIssue)
}

func (ip *IssueProcessor) attachLabels(issue *github.Issue) error {
	labels := open.ParseToGenerateLabels(issue)
	if len(labels) == 0 {
		return nil
	}
	// TODO: check versions to add labels
	// only labels generated do we attach labels to issue
	return ip.Client.AddLabelsToIssue(*(issue.Number), labels)
}

func (ip *IssueProcessor) attachComments(issue *github.Issue) error {
	ip.attachTitleComments(issue)
	ip.attachBodyComments(issue)

	return nil
}

func (ip *IssueProcessor) attachTitleComments(issue *github.Issue) error {
	// check if the title is too short or the body empty.
	if issue.Title != nil && len(*(issue.Title)) > 20 {
		return nil
	}

	// attach comment
	body := fmt.Sprintf(utils.IssueTitleTooShort, *(issue.User.Login), ip.Owner, ip.Repo)
	newComment := &github.IssueComment{
		Body: &body,
	}

	if err := ip.Client.AddCommentToIssue(*(issue.Number), newComment); err != nil {
		return err
	}

	labels := []string{"status/more-info-needed"}
	return ip.Client.AddLabelsToIssue(*(issue.Number), labels)
}

func (ip *IssueProcessor) attachBodyComments(issue *github.Issue) error {
	if issue.Body != nil && len(*(issue.Body)) > 100 {
		return nil
	}

	// attach comment
	body := fmt.Sprintf(utils.IssueDescriptionTooShort, *(issue.User.Login), ip.Owner, ip.Repo, ip.Owner, ip.Repo)
	newComment := &github.IssueComment{
		Body: &body,
	}
	if err := ip.Client.AddCommentToIssue(*(issue.Number), newComment); err != nil {
		return err
	}

	if ip.Client.IssueHasLabel(*(issue.Number), "status/more-info-needed") {
		return nil
	}

	labels := []string{"status/more-info-needed"}
	return ip.Client.AddLabelsToIssue(*(issue.Number), labels)
}

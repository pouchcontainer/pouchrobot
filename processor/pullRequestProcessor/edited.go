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

package pullRequestProcessor

import (
	"fmt"

	"github.com/pouchcontainer/pouchrobot/processor/pullRequestProcessor/open"
	"github.com/pouchcontainer/pouchrobot/utils"

	"github.com/google/go-github/github"
)

// ActToPREdited acts to the event which represents pull request edition.
func (prp *PullRequestProcessor) ActToPREdited(pr *github.PullRequest) error {
	// update labels
	prp.updateLabels(pr)
	// update comment
	prp.updateComments(pr)

	return nil
}

func (prp *PullRequestProcessor) updateLabels(pr *github.PullRequest) error {
	newLabels := open.ParseToGeneratePRLabels(pr)
	if len(newLabels) == 0 {
		return nil
	}

	// get a string slice of labels attached to the current pull request.
	strLabels, err := prp.Client.GetStrLabelsInIssue(*(pr.Number))
	if err != nil {
		return err
	}

	deltaLabels := utils.DeltaSlice(strLabels, newLabels)

	if len(deltaLabels) == 0 {
		return nil
	}

	// add delta labels to pull request
	return prp.Client.AddLabelsToIssue(*(pr.Number), deltaLabels)
}

func (prp *PullRequestProcessor) updateComments(pr *github.PullRequest) error {
	prp.updateTitleComment(pr)
	prp.updateBodyComment(pr)

	return nil
}

func (prp *PullRequestProcessor) updateTitleComment(pr *github.PullRequest) error {
	// check if the title is too short or the body empty.
	if pr.Title == nil || len(*(pr.Title)) < 20 {
		if _, exist := prp.Client.IssueHasComment(*(pr.Number), utils.IssueTitleTooShortSubStr); exist {
			// do nothing
			return nil
		}

		body := fmt.Sprintf(utils.PRTitleTooShort, *(pr.User.Login))
		newComment := &github.IssueComment{
			Body: &body,
		}
		return prp.Client.AddCommentToPR(*(pr.Number), newComment)
	}

	// PR title meets the length
	id, exist := prp.Client.IssueHasComment(*(pr.Number), utils.PRTitleTooShortSubStr)
	if !exist {
		// do nothing
		return nil
	}

	return prp.Client.RemoveComment(id)
}

func (prp *PullRequestProcessor) updateBodyComment(pr *github.PullRequest) error {
	// check if the pull request decription is too short or the body empty.
	if pr.Body == nil || len(*(pr.Body)) < 100 {
		if _, exist := prp.Client.IssueHasComment(*(pr.Number), utils.PRDescriptionTooShortSubStr); exist {
			// do nothing
			return nil
		}

		body := fmt.Sprintf(utils.PRDescriptionTooShort, *(pr.User.Login))
		newComment := &github.IssueComment{
			Body: &body,
		}
		return prp.Client.AddCommentToPR(*(pr.Number), newComment)
	}

	// PR title meets the length
	id, exist := prp.Client.IssueHasComment(*(pr.Number), utils.PRDescriptionTooShortSubStr)
	if !exist {
		// do nothing
		return nil
	}

	return prp.Client.RemoveComment(id)
}

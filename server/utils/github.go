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

package utils

import (
	"encoding/json"

	"github.com/google/go-github/github"
)

// ExtractActionType extracts the action type.
func ExtractActionType(data []byte) (string, error) {
	var m struct {
		Action string `json:"action"`
	}
	if err := json.Unmarshal(data, &m); err != nil {
		return "", err
	}
	return m.Action, nil
}

// ExactIssue extracts the issue from request body.
func ExactIssue(data []byte) (github.Issue, error) {
	var m struct {
		Issue github.Issue `json:"issue"`
	}
	if err := json.Unmarshal(data, &m); err != nil {
		return github.Issue{}, err
	}
	return m.Issue, nil
}

// ExactPR extracts the pull request from request body.
func ExactPR(data []byte) (github.PullRequest, error) {
	var m struct {
		PullRequest github.PullRequest `json:"pull_request"`
	}
	if err := json.Unmarshal(data, &m); err != nil {
		return github.PullRequest{}, err
	}
	return m.PullRequest, nil
}

// ExactIssueComment extracts the issue comment from request body.
func ExactIssueComment(data []byte) (github.IssueComment, error) {
	var m struct {
		IssueComment github.IssueComment `json:"comment"`
	}
	if err := json.Unmarshal(data, &m); err != nil {
		return github.IssueComment{}, err
	}
	return m.IssueComment, nil
}

// ExactIssueLabels extracts the issue labels from request body.
func ExactIssueLabels(data []byte) ([]string, error) {
	var m struct {
		Labels []string `json:"labels"`
	}
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, err
	}
	return m.Labels, nil
}

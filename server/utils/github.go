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

// ExactIssueLables extracts the issue labels from request body.
func ExactIssueLables(data []byte) ([]string, error) {
	var m struct {
		Labels []string `json:"labels"`
	}
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, err
	}
	return m.Labels, nil
}

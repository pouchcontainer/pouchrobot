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

package gh

import (
	"context"
	"strings"

	"github.com/google/go-github/github"
	"github.com/sirupsen/logrus"
)

// ListComments lists all comments in an issue including pull request.
func (c *Client) ListComments(num int) ([]*github.IssueComment, error) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	comments, _, err := c.Client.Issues.ListComments(context.Background(), c.owner, c.repo, num, nil)
	if err != nil {
		logrus.Errorf("failed to list comment in issue(pr) %d: %v", num, err)
		return nil, err
	}
	logrus.Debugf("succeed in listing comments for issue(pr) %d", num)
	return comments, nil
}

// AddCommentToIssue adds comment to an issue.
func (c *Client) AddCommentToIssue(num int, comment *github.IssueComment) error {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	if _, _, err := c.Client.Issues.CreateComment(context.Background(), c.owner, c.repo, num, comment); err != nil {
		logrus.Errorf("failed to add comment %s to issue(pr) %d: %v", *(comment.Body), num, err)
		return err
	}
	logrus.Debugf("succeed in adding comment %s for issue %d", *(comment.Body), num)
	return nil
}

// RemoveComment removes a comment for an issue.
func (c *Client) RemoveComment(id int) error {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	if _, err := c.Client.Issues.DeleteComment(context.Background(), c.owner, c.repo, id); err != nil {
		logrus.Errorf("failed to remove comment %d: %v", id, err)
		return err
	}
	logrus.Debugf("succeed in removing comment %s for issue", id)
	return nil
}

// RmCommentsViaStr removes a comment in an issue via given string
func (c *Client) RmCommentsViaStr(num int, str string) error {
	comments, err := c.ListComments(num)
	if err != nil {
		return err
	}

	for _, comment := range comments {
		if comment.Body != nil && strings.Contains(*(comment.Body), str) {
			return c.RemoveComment(*(comment.ID))
		}
	}
	return nil
}

// RmCommentsViaStrAndAttach removes all comments contains the string str and
// attaches a brand new commnet constructed by body.
// In automan, many cases needs this actions to fresh the comments.
func (c *Client) RmCommentsViaStrAndAttach(num int, str string, body string) error {
	comments, err := c.ListComments(num)
	if err != nil {
		return err
	}

	// remove all the existing CI failure comments
	for _, comment := range comments {
		if comment.Body == nil {
			continue
		}

		if !strings.Contains(*(comment.Body), str) {
			continue
		}

		c.RemoveComment(*(comment.ID))
	}

	newComment := &github.IssueComment{
		Body: &body,
	}

	return c.AddCommentToPR(num, newComment)
}

// IssueHasComment returns true if the issue contains a commnet who has substring of 'elment'
func (c *Client) IssueHasComment(num int, element string) (int, bool) {
	comments, err := c.ListComments(num)
	if err != nil {
		return -1, false
	}
	for _, comment := range comments {
		if comment.Body != nil && strings.Contains(*(comment.Body), element) {
			return *(comment.ID), true
		}
	}
	return -1, false
}

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

	"github.com/google/go-github/github"
	"github.com/sirupsen/logrus"
)

// GetPullRequests gets pull request list for a repo.
func (c *Client) GetPullRequests(opt *github.PullRequestListOptions) ([]*github.PullRequest, error) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	pullRequests, _, err := c.Client.PullRequests.List(context.Background(), c.owner, c.repo, opt)
	if err != nil {
		logrus.Errorf("failed to list pull request in repo %s: %v", c.repo, err)
		return nil, err
	}
	logrus.Debugf("succeed in getting pull requests in repo %s", c.repo)
	return pullRequests, nil
}

// GetSinglePR gets a single PR from repo.
func (c *Client) GetSinglePR(num int) (*github.PullRequest, error) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	pullRequest, _, err := c.Client.PullRequests.Get(context.Background(), c.owner, c.repo, num)
	if err != nil {
		logrus.Errorf("failed to get single pull request %d in repo %s: %v", num, c.repo, err)
		return nil, err
	}
	logrus.Debugf("succeed in getting single pull request %d in repo %s", num, c.repo)
	return pullRequest, nil
}

// ListPRComments lists comments for a pull request.
func (c *Client) ListPRComments(num int) ([]*github.PullRequestComment, error) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	prComments, _, err := c.Client.PullRequests.ListComments(context.Background(), c.owner, c.repo, num, nil)
	if err != nil {
		logrus.Errorf("failed to list comments for pr %d: %v", num, err)
		return nil, err
	}
	logrus.Debugf("succeed in list comments for pr %d:", num)
	return prComments, nil
}

// AddCommentToPR adds comment to a pull request.
func (c *Client) AddCommentToPR(num int, comment *github.IssueComment) error {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	if _, _, err := c.Client.Issues.CreateComment(context.Background(), c.owner, c.repo, num, comment); err != nil {
		logrus.Errorf("failed to add comment %s to pr %d: %v", *(comment.Body), num, err)
		return err
	}
	logrus.Debugf("succeed in creating comment %s for pull request %d", *(comment.Body), num)
	return nil
}

// ListCommits lists all commits in a pull request.
func (c *Client) ListCommits(num int) ([]*github.RepositoryCommit, error) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	commits, _, err := c.PullRequests.ListCommits(context.Background(), c.owner, c.repo, num, nil)
	if err != nil {
		logrus.Errorf("failed to list commits in pull request %d: %v", num, err)
		return nil, err
	}
	logrus.Debugf("succeed in listing commits in pull request %d", num)
	return commits, nil
}

// ListPRReviews lists all reviews on a pull request.
func (c *Client) ListPRReviews(num int) ([]*github.PullRequestReview, error) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	reviews, _, err := c.PullRequests.ListReviews(context.Background(), c.owner, c.repo, num, nil)
	if err != nil {
		logrus.Errorf("failed to list reviews in pull request %d: %v", num, err)
		return nil, err
	}
	logrus.Debugf("succeed in listing reviews in pull request %d", num)
	return reviews, nil
}

// CreatePR creates a brand new pull request in repo.
func (c *Client) CreatePR(newPR *github.NewPullRequest) (*github.PullRequest, error) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	pullRequest, _, err := c.PullRequests.Create(context.Background(), c.owner, c.repo, newPR)
	if err != nil {
		logrus.Errorf("failed to create pull request: %v", err)
		return nil, err
	}
	logrus.Debug("succeed in creating pull request")
	return pullRequest, nil
}

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

	"github.com/pouchcontainer/pouchrobot/utils"

	"github.com/google/go-github/github"
	"github.com/sirupsen/logrus"
)

// GetIssues gets issues of a repo.
func (c *Client) GetIssues(opt *github.IssueListByRepoOptions) ([]*github.Issue, error) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	issues, _, err := c.Client.Issues.ListByRepo(context.Background(), c.owner, c.repo, opt)
	if err != nil {
		logrus.Errorf("failed to list issues in repo %s: %v", c.repo, err)
		return nil, err
	}
	logrus.Debugf("succeed in getting issues in repo %s", c.repo)
	return issues, nil
}

// CreateIssue creates a brand new issue in repo's issue list.
func (c *Client) CreateIssue(title, body string) error {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	issueRequest := &github.IssueRequest{
		Title: &title,
		Body:  &body,
	}
	if _, _, err := c.Issues.Create(context.Background(), c.owner, c.repo, issueRequest); err != nil {
		logrus.Errorf("failed to create issue in repo %s: %v", c.repo, err)
		return err
	}
	return nil
}

// EditIssue modify a specific issue's property in repo's issue list
func (c *Client) EditIssue(opt *github.IssueRequest, order int) (*github.Issue, error) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	issue, _, err := c.Client.Issues.Edit(context.Background(), c.owner, c.repo, order ,opt)
	if err != nil {
		logrus.Errorf("failed to edit issue:%v in repo %s: %v", order, c.repo, err)
		return nil, err
	}
	logrus.Debugf("succeed in edit issue:%v in repo %s", order, c.repo)
	return issue, nil
}

// GetAllLabels gets all labels of a repo, not an issue, nor a pull request
func (c *Client) GetAllLabels() ([]*github.Label, error) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	labels, _, err := c.Client.Issues.ListLabels(context.Background(), c.owner, c.repo, nil)
	if err != nil {
		logrus.Errorf("failed to get labels in repo %s: %v", c.repo, err)
		return nil, err
	}
	logrus.Debugf("succeed in listing all labels in repo %s", c.repo)
	return labels, nil
}

// GetLabelsInIssue gets labels attached on a single issue whose id is num.
func (c *Client) GetLabelsInIssue(num int) ([]*github.Label, error) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	labels, _, err := c.Client.Issues.ListLabelsByIssue(context.Background(), c.owner, c.repo, num, nil)
	if err != nil {
		logrus.Errorf("failed to get labels in issue %d: %v", num, err)
		return nil, err
	}
	logrus.Debugf("succeed in getting labels in issue %d", num)
	return labels, nil
}

// GetStrLabelsInIssue gets string labels attached on a single issue whose id is num.
func (c *Client) GetStrLabelsInIssue(num int) ([]string, error) {
	labels, err := c.GetLabelsInIssue(num)
	if err != nil {
		return nil, err
	}

	result := []string{}
	for _, label := range labels {
		result = append(result, *(label.Name))
	}
	return result, nil
}

// AddLabelsToIssue adds labels to an issue
func (c *Client) AddLabelsToIssue(num int, labels []string) error {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	if _, _, err := c.Client.Issues.AddLabelsToIssue(context.Background(), c.owner, c.repo, num, labels); err != nil {
		logrus.Errorf("failed to add labels %s to issue(pr) %d: %v", labels, num, err)
		return err
	}
	logrus.Debugf("succeed in adding labels %v for issue %d", labels, num)
	return nil
}

// RemoveLabelForIssue removes a label from an issue.
func (c *Client) RemoveLabelForIssue(num int, label string) error {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	if _, err := c.Client.Issues.RemoveLabelForIssue(context.Background(), c.owner, c.repo, num, label); err != nil {
		logrus.Errorf("failed to remove label %s for issue(pr) %d: %v", label, num, err)
		return err
	}
	logrus.Debugf("succeed in removing label %v for issue %d", label, num)
	return nil
}

// ReplaceLabelsForIssue replaces all labels for an issue.
func (c *Client) ReplaceLabelsForIssue(num int, labels []string) error {
	if _, _, err := c.Client.Issues.ReplaceLabelsForIssue(context.Background(), c.owner, c.repo, num, labels); err != nil {
		logrus.Errorf("failed to replace labels %v for issue(pr) %d: %v", labels, num, err)
		return err
	}
	logrus.Debugf("succeed in replacing labels %v for issue %d", labels, num)
	return nil
}

// IssueContainsLabels return whether issue contains labels
func (c *Client) IssueContainsLabels(num int, labels []string) bool {
	rawLabels, err := c.GetLabelsInIssue(num)
	if err != nil {
		return false
	}
	labelSlice := []string{}
	for _, rawLabel := range rawLabels {
		labelSlice = append(labelSlice, *(rawLabel.Name))
	}

	return utils.SliceContainsSlice(labelSlice, labels)
}

// AssignIssueToUsers assigns users to the specified issue.
func (c *Client) AssignIssueToUsers(num int, users []string) error {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	if _, _, err := c.Client.Issues.AddAssignees(context.Background(), c.owner, c.repo, num, users); err != nil {
		logrus.Errorf("failed to assign users %s to issue(pr) %d: %v", users, num, err)
		return err
	}
	logrus.Debugf("succeed in assign users %s for pull request %d", users, num)
	return nil
}

// UnassignIssueToUsers assigns users to the specified issue.
func (c *Client) UnassignIssueToUsers(num int, users []string) error {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	if _, _, err := c.Client.Issues.AddAssignees(context.Background(), c.owner, c.repo, num, users); err != nil {
		logrus.Errorf("failed to assign users %s to issue(pr) %d: %v", users, num, err)
		return err
	}
	logrus.Debugf("succeed in assign users %s for pull request %d", users, num)
	return nil
}

// IssueHasLabel judges if an issue has a specified label.
func (c *Client) IssueHasLabel(num int, inputLabel string) bool {
	labels, err := c.GetLabelsInIssue(num)
	if err != nil {
		return false
	}
	for _, label := range labels {
		if label.GetName() == inputLabel {
			return true
		}
	}
	return false
}

// SearchIssues searches issues.
// search result's wrapper is never be nil.
func (c *Client) SearchIssues(query string, opt *github.SearchOptions, all bool) (*github.IssuesSearchResult, error) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	if all && opt == nil {
		opt = new(github.SearchOptions)
		opt.Page = 1 // first page.
		opt.PerPage = 30
	}

	issueSearchResult := &github.IssuesSearchResult{}

	for {
		result, resp, err := c.Search.Issues(context.Background(), query, opt)
		if err != nil {
			logrus.Errorf("failed to search issues by query %s", query)
			return nil, err
		}
		if result.Total == nil {
			break
		}
		issueSearchResult.Total = result.Total
		issueSearchResult.Issues = append(issueSearchResult.Issues, result.Issues...)

		// just retrieve a page.
		if !all {
			break
		}

		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	return issueSearchResult, nil
}
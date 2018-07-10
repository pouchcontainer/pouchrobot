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

package issueTask

import (
	"time"

	"github.com/google/go-github/github"

	"github.com/sirupsen/logrus"
)

// CloseOutOfDateIssues close out of date issue
func (it *IssueTask) CloseOutOfDateIssues(leftRetryTime int) error {
	logrus.Info("start to close out-of-date issues")

	var err error

	// get all open issue
	opt := &github.IssueListByRepoOptions{
		State: "open",
	}
	issues, e := it.Client.GetIssues(opt)
	if e != nil {
		err = e
	}

	var closedNums int

	// close out-of-day issue
	for _, issue := range issues {
		if it.isOutOfDay(issue) {
			issueID := *(issue.Number)
			logrus.Infof("close issue #%d", issueID)

			e := it.Client.CloseIssue(issueID)

			if e != nil {
				err = e
				logrus.Errorf("failed to close issue #%d: %v", issueID, err)
			} else {
				closedNums++
			}

		}
	}

	// if there is error and there retry time left
	// run closeOutOfDateIssues again
	if err != nil && leftRetryTime > 0 {
		logrus.Warnf("there is error %v retry. there is %d retry time left", err, leftRetryTime)
		it.CloseOutOfDateIssues(leftRetryTime - 1)
	} else {
		logrus.Infof("sucess close the %d out-of-date issues", closedNums)
	}

	return nil
}

// isOutOfDay if updateAt - createAt more than config.MaxDayOfNoActionIssue day then
// the issue is out of day
func (it *IssueTask) isOutOfDay(issue *github.Issue) bool {
	timeSpan := issue.UpdatedAt.Sub(time.Now())
	return int(timeSpan.Hours()) <= 24*it.MaxDayOfNoActionIssue
}

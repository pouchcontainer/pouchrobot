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

package fetcher

import (

	"github.com/google/go-github/github"
	"github.com/sirupsen/logrus"
	"time"
	"context"
)

// CheckPRsConflict checks that if a PR is conflict with the against branch.
func (f *Fetcher) CheckIssue() error {
	logrus.Info("start to check timeout issue")
	opt := &github.IssueListByRepoOptions{
		State: "open",
	}
	prs, err := f.client.GetIssues(opt)
	if err != nil {
		return err
	}

	for _, pr := range prs {
		f.checkIssue(pr)
	}
	return nil
}

func (f *Fetcher) checkIssue(p *github.Issue) error {

	// if PR can be merged to specified branch
	now := time.Now()
	dur, _ := time.ParseDuration("-1s")
	checkLine := now.Add(dur)
	close := "close"
	if p != nil && p.GetUpdatedAt().Before(checkLine){
		// just remove conflict label if there is one
		// and remove conflict comments if there are some
		logrus.Infof("No #%d issue is closed!!", *p.Number)
		issue := github.IssueRequest{
			State : &close,
		}
		_, _, error := f.client.Issues.Edit(context.Background(), "paul-yml", "testrobot", *p.Number, &issue)
		if error != nil {
			logrus.Error("lock error %v", error)
		}
		return nil
	}

	//logrus.Infof("PR %d: found conflict", *(pr.Number))
	//// remove LGTM label if conflict happens
	//if f.client.IssueHasLabel(*(pr.Number), "LGTM") {
	//	f.client.RemoveLabelForIssue(*(pr.Number), "LGTM")
	//}
	//
	//// attach a label and add comments
	//if !f.client.IssueHasLabel(*(pr.Number), utils.PRConflictLabel) {
	//	f.client.AddLabelsToIssue(*(pr.Number), []string{utils.PRConflictLabel})
	//}
	// attach a comment to the pr,
	// and attach a lable confilct/need-rebase to pr

	return nil

}

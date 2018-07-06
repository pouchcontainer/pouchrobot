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
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/pouchcontainer/pouchrobot/utils"

	"github.com/google/go-github/github"
	"github.com/sirupsen/logrus"
)

// CheckPRsGap checks that if a PR is more than 10 commits behind the branch.
func (f *Fetcher) CheckPRsGap() error {
	logrus.Info("start to check PR's gap")
	opt := &github.PullRequestListOptions{
		State: "open",
	}
	prs, err := f.client.GetPullRequests(opt)
	if err != nil {
		return err
	}

	for _, pr := range prs {
		logrus.Info("start to check prs")
		go func(pr *github.PullRequest) {
			if err := f.checkPRGap(pr); err != nil {
				logrus.Errorf("failed to check pull request %d gap: %v", *pr.Number, err)
			}
		}(pr)
	}
	return nil
}

func (f *Fetcher) checkPRGap(p *github.PullRequest) error {
	pr, err := f.client.GetSinglePR(*(p.Number))
	logrus.Infof("start to check pr %d", *(p.Number))
	if err != nil {
		return err
	}

	// get master branch info
	if err := prepareMasterEnv(); err != nil {
		return err
	}
	logrus.Infof("prepare master env done :pr %d", *(p.Number))

	msLogString, err := getLogInfo("master")
	if err != nil {
		return err
	}
	logrus.Infof("get log info done :pr %d", *(p.Number))

	// get pr branch info
	prNum := strconv.Itoa(*p.Number)

	// FIXME handle the situation when failed to clean existing pr branch.
	cleanPrBranchEnv(prNum)
	logrus.Infof("clean exsiting branch done:pr %d", *(p.Number))

	if err = preparePrBranchEnv(prNum); err != nil {
		return err
	}
	logrus.Infof("prepare pr branch env done :pr %d", *(p.Number))

	prBrLogString, err := getLogInfo("pr branch " + prNum)
	if err != nil {
		logrus.Error(err)
		return err
	}
	logrus.Infof("get pr log info done :pr %d", *(p.Number))

	gap := compareAndgetGap(msLogString, prBrLogString)
	logrus.Infof("the gap is %d", gap)
	if gap < 10 {
		return nil
	}

	// continue if gap between master and pr base is over 10
	logrus.Infof("PR %d: found gap %d", *(pr.Number), gap)

	// remove LGTM label if gap happens
	if f.client.IssueHasLabel(*(pr.Number), "LGTM") {
		f.client.RemoveLabelForIssue(*(pr.Number), "LGTM")
	}

	// attach a label and add comments
	if !f.client.IssueHasLabel(*(pr.Number), utils.PRGapLabel) {
		f.client.AddLabelsToIssue(*(pr.Number), []string{utils.PRGapLabel})
	}

	// attach a comment to the pr,
	// and attach a label gap to pr

	return f.AddGapCommentToPR(pr, gap)
}

// AddGapCommentToPR adds gap comments to specific pull request.
func (f *Fetcher) AddGapCommentToPR(pr *github.PullRequest, gap int) error {
	if pr.User == nil || pr.User.Login == nil {
		logrus.Infof("failed to get user from PR %d: empty User", *(pr.Number))
		return nil
	}

	comments, err := f.client.ListComments(*(pr.Number))
	if err != nil {
		return err
	}

	body := fmt.Sprintf(utils.PRGapComment, *(pr.User.Login), strconv.Itoa(gap))
	newComment := &github.IssueComment{
		Body: &body,
	}

	if len(comments) == 0 {
		return f.client.AddCommentToIssue(*(pr.Number), newComment)
	}

	latestComment := comments[len(comments)-1]
	if strings.Contains(*(latestComment.Body), utils.PRGapSubStr) {
		// remove all existing gap comments
		for _, comment := range comments[:(len(comments) - 1)] {
			if strings.Contains(*(comment.Body), utils.PRGapSubStr) {
				if err := f.client.RemoveComment(*(comment.ID)); err != nil {
					continue
				}
			}
		}
		return nil
	}

	// remove all existing gap comments
	for _, comment := range comments {
		if strings.Contains(*(comment.Body), utils.PRGapSubStr) {
			if err := f.client.RemoveComment(*(comment.ID)); err != nil {
				continue
			}
		}
	}

	// add a brand new gap comment
	return f.client.AddCommentToIssue(*(pr.Number), newComment)
}

func prepareMasterEnv() error {
	cmd := exec.Command("git", "checkout", "master")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to checkout master: %v", err)
	}

	cmd = exec.Command("git", "fetch", "upstream", "master")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to git fetch upstreanm master: %v", err)
	}

	cmd = exec.Command("git", "rebase", "upstream/master")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to git rebase upstreanm/master: %v", err)
	}

	return nil
}

func preparePrBranchEnv(prNum string) error {
	cmd := exec.Command("git", "pull", "upstream", "pull/"+prNum+"/head:"+"new-"+prNum)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to pull pr %s: %v", prNum, err)
	}

	cmd = exec.Command("git", "checkout", "new-"+prNum)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to checkout pr %s branch: %v", prNum, err)
	}

	return nil
}

func cleanPrBranchEnv(prNum string) error {
	cmd := exec.Command("git", "branch", "-D", "new-"+prNum)
	// TODO ignore not found error
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to remove existing pr %s branch %v", prNum, err)
	}
	return nil
}

func getLogInfo(branch string) (string, error) {
	var Out bytes.Buffer
	cmd := exec.Command("git", "log", "--oneline")
	cmd.Stdout = &Out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to get %s log: %v", branch, err)
	}

	return Out.String(), nil
}

func compareAndgetGap(msLogString string, prBrLogString string) int {
	var count int

	prBrLog := strings.Split(prBrLogString, "\n")
	msLog := strings.Split(msLogString, "\n")

	for k, v := range prBrLog {
		if strings.Contains(msLogString, v) {
			count = len(prBrLog) - k
			break
		}
	}

	gap := len(msLog) - count
	return gap
}

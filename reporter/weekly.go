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

package reporter

import (
	"fmt"
	"strings"
	"time"

	"github.com/pouchcontainer/pouchrobot/utils"

	"github.com/google/go-github/github"
	"github.com/sirupsen/logrus"
)

// WeekReport contains details about elements to construct a report.
type WeekReport struct {
	// time of this weekly report's start time.
	StartDate string

	// time of this weekly report's end time.
	EndDate string

	// Watch defines currently how many github users are watching this repo.
	Watch int

	// Star defines currently how many github users are staring this repo.
	Star int

	// Fork defines currently how many github users have forked this repo.
	Fork int

	// ContributorsCount defines the number of contributors.
	ContributorsCount int

	// NumOfNewIssues is the issues number which are created in the last week.
	NumOfNewIssues int

	// NumOfClosedIssues is the issues number which are closed in the last week.
	NumOfClosedIssues int

	// MergedPR defines how many pull requests have beem merge between time StartDate and EndDate.
	MergedPR map[string][]*SimplePR

	// CountOfPR defines the number of merged pull request.
	CountOfPR int

	// NewContributors defines new contributors between time StartDate and EndDate.
	NewContributors []string

	// PRReviewsByUser defines that all pull request reviews submitted between time StartDate and EndDate.
	// PRReviewsByUser has a type map, the key is User, Value is the number of pull reuqest reviews of single User.
	PRReviewsByUser map[string]int
}

// SimplePR represents
type SimplePR struct {
	Num     int
	Title   string
	HTMLURL string
}

func (r *Reporter) weeklyReport() error {
	wr, err := r.constructWeekReport()
	if err != nil {
		return err
	}

	issueTitle := fmt.Sprintf("Weekly Report of PouchContainer from %s to %s", wr.StartDate, wr.EndDate)
	issueBody := wr.String()

	return r.client.CreateIssue(issueTitle, issueBody)
}

func (r *Reporter) constructWeekReport() (WeekReport, error) {
	var wr WeekReport

	now := time.Now()
	data := strings.Split(now.String(), " ")
	today := data[0]

	lastWeek := now.Add(-7 * 24 * time.Hour)
	data = strings.Split(lastWeek.String(), " ")
	dayBeforeAWeek := data[0]

	wr.EndDate = today
	wr.StartDate = dayBeforeAWeek

	// get repository details
	repo, err := r.client.GetRepository()
	if err != nil {
		return wr, err
	}

	// get contributor number of repository
	if contributors, err := r.client.ListContributors(); err == nil {
		wr.ContributorsCount = len(contributors)
	}

	wr.Watch = *(repo.SubscribersCount)
	wr.Star = *(repo.StargazersCount)
	wr.Fork = *(repo.ForksCount)

	// get merged pull request details

	logrus.Infof("Start: %s, End: %s", wr.StartDate, wr.EndDate)
	query := fmt.Sprintf("is:merged type:pr repo:%s/%s merged:>=%s", r.client.Owner(), r.client.Repo(), wr.StartDate)
	issueSearchResult, err := r.client.SearchIssues(query, nil, true)
	if err != nil {
		return wr, err
	}

	r.setContributorAndPRSummary(&wr, issueSearchResult)

	r.CalculateReviews(&wr)

	return wr, nil
}

func (r *Reporter) setContributorAndPRSummary(wr *WeekReport, issueSearchResult *github.IssuesSearchResult) {
	wr.CountOfPR = issueSearchResult.GetTotal()
	wr.MergedPR = map[string][]*SimplePR{}

	// SearchIssues returns a list of issue, and we can treat them as pull request as well.
	for _, pr := range issueSearchResult.Issues {
		comments, err := r.client.ListComments(*pr.Number)
		if err != nil {
			continue
		}
		// determine whether this is a new contributor via pull request comments.
		for _, comment := range comments {
			if comment.Body != nil && strings.HasSuffix(*comment.Body, utils.FirstCommitCommentSubStr) {
				wr.NewContributors = append(wr.NewContributors, *pr.User.Login)
				break
			}
		}

		newSimplePR := &SimplePR{
			Title:   *pr.Title,
			HTMLURL: *pr.HTMLURL,
			Num:     *pr.Number,
		}

		if strings.HasPrefix(*pr.Title, "feature:") || strings.HasPrefix(*pr.Title, "feat:") {
			if _, ok := wr.MergedPR["feature"]; !ok {
				wr.MergedPR["feature"] = []*SimplePR{newSimplePR}
			} else {
				wr.MergedPR["feature"] = append(wr.MergedPR["feature"], newSimplePR)
			}
		} else if strings.HasPrefix(*pr.Title, "bugfix:") || strings.HasPrefix(*pr.Title, "fix:") {
			if _, ok := wr.MergedPR["bugfix"]; !ok {
				wr.MergedPR["bugfix"] = []*SimplePR{newSimplePR}
			} else {
				wr.MergedPR["bugfix"] = append(wr.MergedPR["bugfix"], newSimplePR)
			}
		} else if strings.HasPrefix(*pr.Title, "doc:") || strings.HasPrefix(*pr.Title, "docs:") {
			if _, ok := wr.MergedPR["doc"]; !ok {
				wr.MergedPR["doc"] = []*SimplePR{newSimplePR}
			} else {
				wr.MergedPR["doc"] = append(wr.MergedPR["doc"], newSimplePR)
			}
		} else if strings.HasPrefix(*pr.Title, "test:") || strings.HasPrefix(*pr.Title, "tests:") {
			if _, ok := wr.MergedPR["test"]; !ok {
				wr.MergedPR["test"] = []*SimplePR{newSimplePR}
			} else {
				wr.MergedPR["test"] = append(wr.MergedPR["test"], newSimplePR)
			}
		} else if _, ok := wr.MergedPR["others"]; !ok {
			wr.MergedPR["others"] = []*SimplePR{newSimplePR}
		} else {
			wr.MergedPR["others"] = append(wr.MergedPR["others"], newSimplePR)
		}
	}

	// make contributor name unique in weekly report.
	wr.NewContributors = utils.UniqueElementSlice(wr.NewContributors)

	return
}

// String returns a string of Week Report
func (wr *WeekReport) String() string {
	totalStr := fmt.Sprintf(`
# Weekly Report of PouchContainer
This is a weekly report of PouchContainer. It summarizes what have changed in PouchContainer in the passed week, including pr merged, new contributors, and more things in the future. 
It is all done by @pouchrobot which is an AI robot.  See: https://github.com/pouchcontainer/pouchrobot.
`)

	// get repo update for this week
	repoUpdateContent := wr.getRepoUpdateContent()
	totalStr += repoUpdateContent

	// get repo update for this week
	prUpdateContent := wr.getPRUpdateContent()
	totalStr += prUpdateContent

	// construct code review details of the past week
	prReviewContent := wr.getPRReviewContent()
	totalStr += prReviewContent

	// calculate new contributors of this week.
	newContributorsContent := wr.getNewContributorsContent()
	totalStr += newContributorsContent

	return totalStr
}

func (wr *WeekReport) getRepoUpdateContent() string {
	header := "## Repo Update \n"

	foreword := ""

	repoUpdate := `
| Watch | Star | Fork | Contributors | New Issues | Closed Issues |
|:-----:|:----:|:----:|:------------:|:----------:|:-------------:|
`
	repoUpdate += fmt.Sprintf("|%d|%d|%d|%d|%d|%d|\n\n", wr.Watch, wr.Star, wr.Fork, wr.ContributorsCount, wr.NumOfNewIssues, wr.NumOfClosedIssues)

	wholeContent := header + foreword + repoUpdate
	return wholeContent
}

func (wr *WeekReport) getPRUpdateContent() string {
	header := fmt.Sprintf(`
		## PR Update
		
		Thanks to contributions from community, PouchContainer team merged %d pull requests in the PouchContainer repositories last week. All these pull requests could be divided into **feature**, **bugfix**, **doc**, **test** and **others**:
		
		`, wr.CountOfPR)

	foreword := ""

	prUpdateContent := ""
	for _, typeStr := range []string{"feature", "bugfix", "doc", "test", "others"} {
		if len(wr.MergedPR[typeStr]) == 0 {
			// if no this type pr merged, no related thing output.
			continue
		}

		var appendStr string
		if typeStr == "feature" {
			appendStr = fmt.Sprintf("### %s 🆕 🔫 \n\n", typeStr)
		} else if typeStr == "bugfix" {
			appendStr = fmt.Sprintf("### %s 🐛 🔪 \n\n", typeStr)
		} else if typeStr == "doc" {
			appendStr = fmt.Sprintf("### %s 📜 📝 \n\n", typeStr)
		} else if typeStr == "test" {
			appendStr = fmt.Sprintf("### %s ✅ ☑️ \n\n", typeStr)
		} else {
			appendStr = fmt.Sprintf("### %s\n\n", typeStr)
		}

		prUpdateContent += appendStr
		for _, pr := range wr.MergedPR[typeStr] {
			prUpdateContent += fmt.Sprintf("* %s ([#%d](%s))\n", pr.Title, pr.Num, pr.HTMLURL)
		}
		prUpdateContent += "\n"
	}

	wholeContent := header + foreword + prUpdateContent
	return wholeContent
}

func (wr *WeekReport) getPRReviewContent() string {
	header := "## Code Review Statistics 🐞 🐞 🐞 \n"

	foreword := "PouchContainer encourages everyone to participant in code review, in order to improve software quality. Everyweek @pouchrobot would automatically help to count pull request reviews of single github user as the following. So, try to help review code in this project.\n\n"

	tableHeader := `| Contributor ID | Pull Request Reviews |
|:--------: | :--------:|
`

	tableContent := ""

	// sort the users
	length := len(wr.PRReviewsByUser)
	users := make([]string, 0, length)
	reviewNums := make([]int, 0, length)
	for user, num := range wr.PRReviewsByUser {
		users = append(users, user)
		reviewNums = append(reviewNums, num)
	}
	for i := 0; i < length-1; i++ {
		for j := i + 1; j < length; j++ {
			if reviewNums[i] < reviewNums[j] {
				reviewNums[i], reviewNums[j] = reviewNums[j], reviewNums[i]
				users[i], users[j] = users[j], users[i]
			}
		}
	}

	// after sorting, construct table content via sorted data
	for i := 0; i < length; i++ {
		tableRow := fmt.Sprintf("|@%s|%d|\n", users[i], reviewNums[i])
		tableContent += tableRow
	}

	tableContent += "\n\n"

	wholeContent := header + foreword + tableHeader + tableContent
	return wholeContent
}

func (wr *WeekReport) getNewContributorsContent() string {
	header := "## New Contributors 🎖 🎖 🎖 \n\n"

	newContributorsContent := ""
	if len(wr.NewContributors) != 0 {
		newContributorsContent += `It is PouchContainer team's great honor to have new contributors in Pouch's community. We really appreciate your contributions. Feel free to tell us if you have any opinion and please share PouchContainer with more people if you could. If you hope to be a contributor as well, please start from https://github.com/alibaba/pouch/blob/master/CONTRIBUTING.md . 🎁 👏 🍺
Here is the list of new contributors:
`
		for _, contributor := range wr.NewContributors {
			newContributorsContent += fmt.Sprintf("@%s\n", contributor)
		}
	} else {
		newContributorsContent += `We have no new contributors in PouchContainer project this week.
PouchContainer team encourages everything about contribution from community.
For more details, please refer to https://github.com/alibaba/pouch/blob/master/CONTRIBUTING.md . 🍻
`
	}

	newContributorsContent += fmt.Sprintf("\n\n Thank all of you!")

	wholeContent := header + newContributorsContent
	return wholeContent
}

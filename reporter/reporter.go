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

	"github.com/google/go-github/github"
	"github.com/pouchcontainer/pouchrobot/gh"
	"github.com/pouchcontainer/pouchrobot/utils"

	"github.com/sirupsen/logrus"
)

// Reporter is a reporter to report weekly update on Github Repo in issues.
// Reports needs to use the repo name to construct the weeklyreport name.
// Since weekly report generation should be general to all kinds of open source project.
type Reporter struct {
	client *gh.Client

	// ReportDay representing which is the weekly report generation day.
	ReportDay string

	// ReportHour representing which is the weekly report generation time on ReportDay.
	ReportHour int

	// Owner is the organization of open source project.
	owner string

	// Repo is the repository name.
	repo string
}

var statsLastWeek = &StatsLastWeek{}

// New initializes a brand new reporter.
func New(client *gh.Client, day string, hour int) *Reporter {
	return &Reporter{
		client:     client,
		owner:      client.Owner(),
		repo:       client.Repo(),
		ReportDay:  day,
		ReportHour: hour,
	}
}

// Run starts to work on reporting things for repo.
func (r *Reporter) Run() {
	logrus.Infof("start to run reporter")

	// initialize fork, star and watch informations
	go r.initRepoInfo(statsLastWeek)

	// Wait time goes to Friday.
	for {
		if time.Now().Weekday().String() == r.ReportDay {
			hour, _, _ := time.Now().Clock()
			logrus.Infof("now time is %s:%d", "Friday", hour)
			if hour == r.ReportHour {
				break
			}
		}
		time.Sleep(30 * time.Minute)
	}

	for {
		// only fixed day, code will enter this for loop block.
		go r.weeklyReport()

		// report one issue every week.
		time.Sleep(7 * 24 * time.Hour)
	}
}

func (r *Reporter) weeklyReport() error {
	logrus.Infof("weekly report generation is truly starting.....")
	// first, construct weekly report data via fresh data
	wr, err := r.constructWeekReport()
	if err != nil {
		return err
	}

	// start to post an issue which represents the weekly report, like:
	// https://github.com/alibaba/pouch/issues/2067
	issueTitle := fmt.Sprintf("WeeklyReport of %s from %s to %s", wr.repo, wr.StartDate, wr.EndDate)
	issueBody := wr.String()

	// after constructing weekly report, use this week data to replace stats last week
	statsLastWeek.Contributors = wr.Contributors
	statsLastWeek.Star = wr.Star
	statsLastWeek.Fork = wr.Fork
	statsLastWeek.Watch = wr.Watch

	return r.client.CreateIssue(issueTitle, issueBody)
}

func (r *Reporter) constructWeekReport() (WeekReport, error) {
	var wr WeekReport

	wr.owner = r.owner
	// the following block is specified for PouchContainer.
	if r.repo == "pouch" {
		r.repo = "PouchContainer"
	}
	wr.repo = r.repo

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

	// list contributors of repository
	listContributorsOpt := &github.ListContributorsOptions{
		ListOptions: github.ListOptions{
			PerPage: 120,
		},
	}
	if contributors, err := r.client.ListContributors(listContributorsOpt); err == nil {
		wr.Contributors = len(contributors)
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
}

// initRepoInfo gets repo's status at the start time of robot.
// This will leads to inaccuracy of first week.
// But week after first one will be correct.
func (r *Reporter) initRepoInfo(lastweek *StatsLastWeek) {
	// get repository details
	repo, err := r.client.GetRepository()
	if err != nil {
		return
	}

	// list contributors of repository
	listContributorsOpt := &github.ListContributorsOptions{
		ListOptions: github.ListOptions{
			PerPage: 120,
		},
	}
	if contributors, err := r.client.ListContributors(listContributorsOpt); err == nil {
		lastweek.Contributors = len(contributors)
	}

	lastweek.Watch = *(repo.SubscribersCount)
	lastweek.Star = *(repo.StargazersCount)
	lastweek.Fork = *(repo.ForksCount)
}

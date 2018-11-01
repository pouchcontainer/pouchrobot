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

import "fmt"

// WeekReport contains details about elements to construct a report.
type WeekReport struct {
	// Owner is the organization of open source project.
	owner string

	// Repo is the repository name.
	repo string

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

	// Contributors defines the number of contributors.
	Contributors int

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

// StatsLastWeek collects repo data from last week.
type StatsLastWeek struct {
	Watch        int
	Star         int
	Fork         int
	Contributors int
}

// SimplePR represents
type SimplePR struct {
	Num     int
	Title   string
	HTMLURL string
}

// String returns a string of Week Report
func (wr *WeekReport) String() string {
	totalStr := fmt.Sprintf(`
# Weekly Report of %s

This is a weekly report of %s. It summarizes what have changed in the project during the passed week, including pr merged, new contributors, and more things in the future. 
It is all done by @pouchrobot which is an AI robot.  See: https://github.com/pouchcontainer/pouchrobot.
`, wr.repo, wr.repo)

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
	repoUpdate += fmt.Sprintf("|%d (↑%d)|%d (↑%d)|%d (↑%d)|%d (↑%d)|%d|%d|\n\n",
		wr.Watch, wr.Watch-statsLastWeek.Watch,
		wr.Star, wr.Star-statsLastWeek.Star,
		wr.Fork, wr.Fork-statsLastWeek.Fork,
		wr.Contributors, wr.Contributors-statsLastWeek.Contributors,
		wr.NumOfNewIssues, wr.NumOfClosedIssues)

	wholeContent := header + foreword + repoUpdate
	return wholeContent
}

func (wr *WeekReport) getPRUpdateContent() string {
	header := fmt.Sprintf(`
## PR Update
		
Thanks to contributions from community, %s team merged **%d** pull requests in the repository last week. All these pull requests could be divided into **feature**, **bugfix**, **doc**, **test** and **others**:
		
`, wr.repo, wr.CountOfPR)

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

	foreword := "This project encourages everyone to participant in code review, in order to improve software quality. Every week @pouchrobot would automatically help to count pull request reviews of single github user as the following. So, try to help review code in this project.\n\n"

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
		newContributorsContent += fmt.Sprintf(`It is %s team's great honor to have new contributors from community. We really appreciate your contributions. Feel free to tell us if you have any opinion and please share this open source project with more people if you could. If you hope to be a contributor as well, please start from https://github.com/alibaba/pouch/blob/master/CONTRIBUTING.md . 🎁 👏 🍺
Here is the list of new contributors:

`, wr.repo)
		for _, contributor := range wr.NewContributors {
			newContributorsContent += fmt.Sprintf("@%s\n", contributor)
		}
	} else {
		newContributorsContent += fmt.Sprintf(`We have no new contributors in this project this week.
%s team encourages everything about contribution from community.
For more details, please refer to https://github.com/alibaba/pouch/blob/master/CONTRIBUTING.md . 🍻
`, wr.repo)
	}

	newContributorsContent += fmt.Sprintf("\n\n Thank all of you!")

	wholeContent := header + newContributorsContent
	return wholeContent
}

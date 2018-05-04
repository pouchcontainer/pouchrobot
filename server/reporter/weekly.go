package reporter

import (
	"fmt"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/google/go-github/github"
	"github.com/pouchcontainer/pouchrobot/server/utils"
)

// WeekReport contains details about elements to construct a report.
type WeekReport struct {
	StartDate         string
	EndDate           string
	Watch             int
	Star              int
	Fork              int
	ContributorsCount int
	MergedPR          map[string][]*SimplePR
	CountOfPR         int
	NewContributors   []string
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

	issueTitle := fmt.Sprintf("Weekly Report of Pouch from %s to %s", wr.StartDate, wr.EndDate)
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
# Weekly Report of Pouch

This is a weekly report of Pouch. It summarizes what have changed in Pouch in the past week, including pull requests which are merged, new contributors, and more things in the future. 
It is all done by @pouchrobot which is an AI robot.

## Repo Update 

| Type      |    Count |
| :-------- | --------:| 
| Watch     |   %d |  
| Star      |   %d |  
| Fork      |   %d |
`,
		wr.Watch,
		wr.Star,
		wr.Fork,
	)

	prUpdateSubStr := fmt.Sprintf(`
## PR Update

Thanks to contributions from community, Pouch team merged %d pull requests in the Pouch repositories last week. All these pull requests could be divided into **feature**, **bugfix**, **doc**, **test** and **others**:

`,
		wr.CountOfPR,
	)
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
		prUpdateSubStr = prUpdateSubStr + appendStr
		for _, pr := range wr.MergedPR[typeStr] {
			prUpdateSubStr = prUpdateSubStr + fmt.Sprintf("* %s ([#%d](%s))\n", pr.Title, pr.Num, pr.HTMLURL)
		}
		prUpdateSubStr = prUpdateSubStr + "\n"
	}
	totalStr = totalStr + prUpdateSubStr

	// calculate new contributors of this week.
	newContribSubstr := "## New Contributors 🎖 🎖 🎖 \n\n"
	if len(wr.NewContributors) != 0 {
		newContribSubstr = newContribSubstr + `It is Pouch team's great honor to have new contributors in Pouch's community. We really appreciate your contributions. Feel free to tell us if you have any opinion and please share Pouch with more people if you could. If you hope to be a contributor as well, please start from https://github.com/alibaba/pouch/blob/master/CONTRIBUTING.md . 🎁 👏 🍺 

Here is the list of new contributors:

`
		for _, contributor := range wr.NewContributors {
			newContribSubstr = newContribSubstr + fmt.Sprintf("@%s\n", contributor)
		}
	} else {
		newContribSubstr = newContribSubstr + `We have no new contributors in Pouch project this week.
Pouch team encourages everything about contribution from community.
For more details, please refer to https://github.com/alibaba/pouch/blob/master/CONTRIBUTING.md . 🍻 
`
	}

	newContribSubstr = newContribSubstr + fmt.Sprintf("\n\n Thank all of you!")
	totalStr = totalStr + newContribSubstr

	return totalStr
}

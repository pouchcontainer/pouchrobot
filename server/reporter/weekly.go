package reporter

import (
	"fmt"
	"strings"

	"github.com/allencloud/automan/server/utils"
	"github.com/google/go-github/github"
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
	Num   int
	Title string
	URL   string
}

func (r *Reporter) weeklyReport() error {
	wr, err := r.construcWeekReport()
	if err != nil {
		return err
	}

	issueTitle := fmt.Sprintf("Weekly Report in Pouch %s - %s", wr.StartDate, wr.EndDate)
	issueBody := wr.String()

	return r.client.CreateIssue(issueTitle, issueBody)
}

func (r *Reporter) construcWeekReport() (WeekReport, error) {
	var wr WeekReport

	// get repository details
	repo, err := r.client.GetRepository()
	if err != nil {
		return wr, err
	}

	wr.Watch = *(repo.WatchersCount)
	wr.Star = *(repo.StargazersCount)
	wr.Fork = *(repo.ForksCount)

	// get merged pull request details
	query := "is:merged type:pr repo:moby/moby merged:>=2017-11-23"
	issueSearchResult, err := r.client.SearchIssues(query, nil)
	if err != nil {
		return wr, err
	}

	// SearchIssues returns a list of issue, and we can treat them as pull request as well.
	prs := issueSearchResult.Issues

	wr.setContributorAndCommits(prs)

	return wr, nil
}

func (wr *WeekReport) setContributorAndCommits(prs []github.Issue) {
	wr.CountOfPR = len(prs)
	for _, pr := range prs {
		if pr.Body != nil && strings.HasSuffix(*pr.Body, utils.FirstCommitCommentSubStr) {
			wr.NewContributors = append(wr.NewContributors, *pr.User.Login)
		}

		newSimplePR := &SimplePR{
			Title: *pr.Title,
			URL:   *pr.URL,
			Num:   *pr.Number,
		}

		if strings.HasPrefix(*pr.Title, "feature:") || strings.HasPrefix(*pr.Title, "feat:") {
			if wr.MergedPR["feature"] == nil {
				wr.MergedPR["feature"] = []*SimplePR{newSimplePR}
			} else {
				wr.MergedPR["feature"] = append(wr.MergedPR["feature"], newSimplePR)
			}
		} else if strings.HasPrefix(*pr.Title, "bugfix:") || strings.HasPrefix(*pr.Title, "fix:") {
			if wr.MergedPR["bugfix"] == nil {
				wr.MergedPR["bugfix"] = []*SimplePR{newSimplePR}
			} else {
				wr.MergedPR["bugfix"] = append(wr.MergedPR["bugfix"], newSimplePR)
			}
		} else if strings.HasPrefix(*pr.Title, "doc:") || strings.HasPrefix(*pr.Title, "docs:") {
			if wr.MergedPR["doc"] == nil {
				wr.MergedPR["doc"] = []*SimplePR{newSimplePR}
			} else {
				wr.MergedPR["doc"] = append(wr.MergedPR["doc"], newSimplePR)
			}
		} else if strings.HasPrefix(*pr.Title, "test:") || strings.HasPrefix(*pr.Title, "tests:") {
			if wr.MergedPR["test"] == nil {
				wr.MergedPR["test"] = []*SimplePR{newSimplePR}
			} else {
				wr.MergedPR["test"] = append(wr.MergedPR["test"], newSimplePR)
			}
		} else {
			if wr.MergedPR["others"] == nil {
				wr.MergedPR["others"] = []*SimplePR{newSimplePR}
			} else {
				wr.MergedPR["others"] = append(wr.MergedPR["others"], newSimplePR)
			}
		}
	}
	return
}

// String returns a string of Week Report
func (wr *WeekReport) String() string {
	totalStr := fmt.Sprintf(`
# Weekly Report in Pouch

%s - %s

## Repo Update 

| Type      |    Count |
| :-------- | --------:| 
| Watch     |   %d |  
| Star      |   %d |  
| Fork      |   %d | 
`,
		wr.StartDate,
		wr.EndDate,
		wr.Watch,
		wr.Star,
		wr.Fork,
	)

	prUpdateSubStr := fmt.Sprintf("## PR Update\n\nLast week, we merged %d pull requests in the Pouch repositories.\n\n", wr.CountOfPR)
	for _, typeStr := range []string{"feature", "bugfix", "doc", "test", "others"} {
		prUpdateSubStr = prUpdateSubStr + fmt.Sprintf("### %s\n\n", typeStr)
		for _, pr := range wr.MergedPR[typeStr] {
			prUpdateSubStr = prUpdateSubStr + fmt.Sprintf("* [%s](%s)\n", pr.Title, pr.URL)
		}
		prUpdateSubStr = prUpdateSubStr + "\n"
	}
	totalStr = totalStr + prUpdateSubStr

	// calculate new contributors of this week.
	newContribSubstr := "## New Contributors\n\n"
	for _, contributor := range wr.NewContributors {
		newContribSubstr = newContribSubstr + fmt.Sprintf("@%s\n", contributor)
	}
	newContribSubstr = newContribSubstr + fmt.Sprintf("\n")
	totalStr = totalStr + newContribSubstr

	return totalStr
}

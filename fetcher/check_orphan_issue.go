package fetcher

import (
	"github.com/sirupsen/logrus"
	"github.com/google/go-github/github"
	"time"
)

// checkOrphanIssue checks that if some issues have no comments and nobody care
func (f *Fetcher) CheckOrphanIssue() error {
	logrus.Info("start to check orphan issues...")

	now :=time.Now()
	// TODO: configurable
	var thirtyDayAgo = now.AddDate(0,0,-30)
	opt := &github.IssueListByRepoOptions{
		State: "open",
		Sort: "comments",
		Direction:"asc",
		Since: thirtyDayAgo,
	}
	isues,err := f.client.GetIssues(opt)
	if err != nil {
		return err
	}

	// stop if all have comments
	if len(isues) == 0 || *(isues[0].Comments) > 0{
		logrus.Info("no orphan issue found...")
		return nil
	}

	// process orphan issues
	for _, isue := range isues {
		f.closeOrphanIssue(isue)
	}

	return nil
}

// close specific Issue
func (f *Fetcher) closeOrphanIssue(isue *github.Issue) error {
	// check comments
	//var err error // why warn??
	if *(isue.Comments) == 0 {
		s := "closed"
		parm := &github.IssueRequest{
			State: &s,
		}
		_, err := f.client.EditIssue(parm, *(isue.ID))
		return err
	}
	return nil
}

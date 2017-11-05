package fetcher

import (
	"fmt"
	"strings"
	"time"

	"github.com/allencloud/automan/server/gh"
	"github.com/allencloud/automan/server/utils"
	"github.com/google/go-github/github"

	"github.com/sirupsen/logrus"
)

// FETCHINTERVAL refers the interval of fetch action
const FETCHINTERVAL = 1 * time.Minute

// Fetcher is a worker to periodically get elements from github
type Fetcher struct {
	client *gh.Client
}

// NewFetcher creates
func NewFetcher(client *gh.Client) *Fetcher {
	return &Fetcher{
		client: client,
	}
}

// Work starts periodical work
func (f *Fetcher) Work() {
	for {
		time.Sleep(FETCHINTERVAL)

		f.CheckPRsConflict()
	}
}

// CheckPRsConflict checks that if a PR is conflict with the against branch.
func (f *Fetcher) CheckPRsConflict() error {
	logrus.Info("start to check PR's conflict")
	opt := &github.PullRequestListOptions{
		State: "open",
	}
	prs, err := f.client.GetPullRequests(opt)
	if err != nil {
		return err
	}

	for _, one := range prs {
		pr, err := f.client.GetSinglePR(*(one.Number))
		if err != nil {
			continue
		}
		if pr.Mergeable != nil && *(pr.Mergeable) == false {
			logrus.Infof("found pull request %d conflict", *(pr.Number))
			// attach a comment to the pr,
			// and attach a lable confilct/need-rebase to pr
			if !hasConflictLabel(f.client, pr) {
				f.AddConflictLabelToPR(pr)
			}
			if !hasConflictComment(f.client, pr) {
				f.AddConflictCommentToPR(pr)
			}
		} else if pr.Mergeable != nil && *(pr.Mergeable) == true {
			if hasConflictLabel(f.client, pr) {
				f.client.RemoveLabelForIssue(*(pr.Number), "conflict/needs-rebase")
			}
		}
	}
	return nil
}

func hasConflictLabel(c *gh.Client, pr *github.PullRequest) bool {
	labels, err := c.GetLabelsInIssue(*(pr.Number))
	if err != nil {
		return false
	}

	for _, label := range labels {
		if *(label.Name) == utils.ConflictLabel {
			return true
		}
	}
	return false
}

func hasConflictComment(c *gh.Client, pr *github.PullRequest) bool {
	comments, err := c.ListComments(*(pr.Number))
	if err != nil {
		return false
	}

	for _, comment := range comments {
		logrus.Infof("pull request %d has comment %s", *(pr.Number), *(comment.Body))
		if strings.Contains(*(comment.Body), "Conflict happens after merging a previous commit. Please rebase the branch against master and push it back again. Thanks a lot.") {
			return true
		}
	}
	return false
}

// AddConflictCommentToPR adds conflict comments to specific pull request.
func (f *Fetcher) AddConflictCommentToPR(pr *github.PullRequest) error {
	newComment := &github.IssueComment{}
	if pr.User == nil || pr.User.Login == nil {
		logrus.Infof("failed to get user from PR %d: empty User", *(pr.Number))
		return nil
	}
	body := fmt.Sprintf(utils.PRConflictComment, *(pr.User.Login))
	newComment.Body = &body
	return f.client.AddCommentToIssue(*(pr.Number), newComment)
}

// AddConflictLabelToPR adds a label of conflict/need-rebase for pull request.
func (f *Fetcher) AddConflictLabelToPR(pr *github.PullRequest) error {
	labels := []string{utils.ConflictLabel}
	return f.client.AddLabelsToIssue(*(pr.Number), labels)
}

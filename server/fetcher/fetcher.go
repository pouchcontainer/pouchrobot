package fetcher

import (
	"fmt"
	"time"

	"github.com/allencloud/automan/server/gh"
	putils "github.com/allencloud/automan/server/processor/utils"
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

	for _, pr := range prs {
		if pr.Mergeable != nil && *(pr.Mergeable) == false {
			logrus.Infof("found pull request %d conflict", *(pr.Number))
			// attach a comment to the pr,
			// and attach a lable confilct/need-rebase to pr
			f.AddConflictCommentToPR(pr)
			f.AddConflictLabelToPR(pr)
		}
	}
	return nil
}

// AddConflictCommentToPR adds conflict comments to specific pull request.
func (f *Fetcher) AddConflictCommentToPR(pr *github.PullRequest) error {
	newComment := &github.IssueComment{}
	if pr.User == nil || pr.User.Login == nil {
		logrus.Infof("failed to get user from PR %d: empty User", *(pr.Number))
		return nil
	}
	body := fmt.Sprintf(putils.PRConflictComment, *(pr.User.Login))
	newComment.Body = &body
	return f.client.AddCommentToIssue(*(pr.Number), newComment)
}

// AddConflictLabelToPR adds a label of conflict/need-rebase for pull request.
func (f *Fetcher) AddConflictLabelToPR(pr *github.PullRequest) error {
	labels := []string{"conflict/need-rebase"}
	return f.client.AddLabelsToIssue(*(pr.Number), labels)
}

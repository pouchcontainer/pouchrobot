package fetcher

import (
	"fmt"
	"strings"

	"github.com/allencloud/automan/server/utils"
	"github.com/google/go-github/github"
	"github.com/sirupsen/logrus"
)

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
		f.checkPRConflict(pr)
	}
	return nil
}

func (f *Fetcher) checkPRConflict(p *github.PullRequest) error {
	pr, err := f.client.GetSinglePR(*(p.Number))
	if err != nil {
		return nil
	}

	if pr.Mergeable == nil || *(pr.Mergeable) == true {
		f.client.RemoveLabelForIssue(*(pr.Number), utils.PRConflictLabel)
		f.client.RemoveCommentViaString(*(pr.Number), utils.PRConflictSubStr)
		return nil
	}

	logrus.Infof("PR %d: found conflict", *(pr.Number))
	if f.client.IssueHasLabel(*(pr.Number), utils.PRConflictLabel) {
		return nil
	}
	// attach a comment to the pr,
	// and attach a lable confilct/need-rebase to pr
	f.client.AddLabelsToIssue(*(pr.Number), []string{utils.PRConflictLabel})
	f.AddConflictCommentToPR(pr)
	return nil
}

// AddConflictCommentToPR adds conflict comments to specific pull request.
func (f *Fetcher) AddConflictCommentToPR(pr *github.PullRequest) error {
	comments, err := f.client.ListComments(*(pr.Number))
	if err != nil {
		return err
	}
	logrus.Infof("PR %d: There are %d comments", *(pr.Number), len(comments))

	if len(comments) == 0 {
		return nil
	}
	latestComment := comments[len(comments)-1]
	if strings.Contains(*(latestComment.Body), utils.PRConflictSubStr) {
		logrus.Infof("PR %d: latest comment %s \nhas\n %s", *(pr.Number), *(latestComment.Body), utils.PRConflictSubStr)
		// do nothing
		return nil
	}

	// remove all existing conflict comments
	for _, comment := range comments {
		if strings.Contains(*(comment.Body), utils.PRConflictSubStr) {
			if err := f.client.RemoveComment(*(comment.ID)); err != nil {
				continue
			}
		}
	}

	// add a brand new conflict comment
	newComment := &github.IssueComment{}
	if pr.User == nil || pr.User.Login == nil {
		logrus.Infof("failed to get user from PR %d: empty User", *(pr.Number))
		return nil
	}
	body := fmt.Sprintf(utils.PRConflictComment, *(pr.User.Login))
	newComment.Body = &body
	return f.client.AddCommentToIssue(*(pr.Number), newComment)
}

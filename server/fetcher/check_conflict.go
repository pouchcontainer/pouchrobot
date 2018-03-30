package fetcher

import (
	"fmt"
	"strings"

	"github.com/pouchcontainer/pouchrobot/server/utils"

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

	// if PR can be merged to specified branch
	if pr.Mergeable == nil || *(pr.Mergeable) == true {
		// just remove conflict label if there is one
		// and remove conflict comments if there are some
		if f.client.IssueHasLabel(*(pr.Number), utils.PRConflictLabel) {
			f.client.RemoveLabelForIssue(*(pr.Number), utils.PRConflictLabel)
		}
		f.client.RmCommentsViaStr(*(pr.Number), utils.PRConflictSubStr)
		return nil
	}

	logrus.Infof("PR %d: found conflict", *(pr.Number))
	// remove LGTM label if conflict happens
	if f.client.IssueHasLabel(*(pr.Number), "LGTM") {
		f.client.RemoveLabelForIssue(*(pr.Number), "LGTM")
	}

	// attach a label and add comments
	if !f.client.IssueHasLabel(*(pr.Number), utils.PRConflictLabel) {
		f.client.AddLabelsToIssue(*(pr.Number), []string{utils.PRConflictLabel})
	}
	// attach a comment to the pr,
	// and attach a lable confilct/need-rebase to pr

	return f.AddConflictCommentToPR(pr)

}

// AddConflictCommentToPR adds conflict comments to specific pull request.
func (f *Fetcher) AddConflictCommentToPR(pr *github.PullRequest) error {
	if pr.User == nil || pr.User.Login == nil {
		logrus.Infof("failed to get user from PR %d: empty User", *(pr.Number))
		return nil
	}

	comments, err := f.client.ListComments(*(pr.Number))
	if err != nil {
		return err
	}
	//logrus.Infof("PR %d: There are %d comments", *(pr.Number), len(comments))

	body := fmt.Sprintf(utils.PRConflictComment, *(pr.User.Login))
	newComment := &github.IssueComment{
		Body: &body,
	}

	if len(comments) == 0 {
		return f.client.AddCommentToIssue(*(pr.Number), newComment)
	}

	latestComment := comments[len(comments)-1]
	if strings.Contains(*(latestComment.Body), utils.PRConflictSubStr) {
		//logrus.Infof("PR %d: latest comment %s \nhas\n %s", *(pr.Number), *(latestComment.Body), utils.PRConflictSubStr)
		// remove all existing conflict comments
		for _, comment := range comments[:(len(comments) - 1)] {
			if strings.Contains(*(comment.Body), utils.PRConflictSubStr) {
				if err := f.client.RemoveComment(*(comment.ID)); err != nil {
					continue
				}
			}
		}
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
	return f.client.AddCommentToIssue(*(pr.Number), newComment)
}

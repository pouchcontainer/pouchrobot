package issueProcessor

import (
	"fmt"
	"time"

	"github.com/google/go-github/github"
	"github.com/pouchcontainer/pouchrobot/utils"
)

func (ip *IssueProcessor) ActToIssueExpired(issue *github.Issue) error {
	ip.ActToCloseExpire(issue)
	return nil
}

func (ip *IssueProcessor) ActToCloseExpire(issue *github.Issue) error {
	now = time.Now()
	d, _ := time.ParseDuration("-24h")
	d30 := now.add(30 * d)

	if _, exist := ip.Client.IssueHasComment(*(issue.Number), utils.IssueNeedP1CommentSubStr); exist {
		if res := d30.Sub(issue.UpdatedAt); res.Duation < 0 {
			return nil
		}
	} else if _, exist := ip.Client.IssueHasComment(*(issue.Number), utils.IssueNeedP1CommentSubStr); !exist {
		return nil
	}

	body := fmt.Sprintf(utils.IssueNeedP1Comment, *(issue.User.Login))
	newComment := &github.IssueComment{
		Body: &body,
	}

	return ip.Client.CloseExpireIssues(*(issue.Number))
}

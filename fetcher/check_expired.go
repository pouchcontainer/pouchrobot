package fetcher

import (
	"github.com/sirupsen/logrus"
	"github.com/google/go-github/github"
	"github.com/alibaba/pouch/pkg/utils"
	"time"
	"strconv"
	"net/http"
	"bytes"
)

const EXPIRETIME = 1 * utils.Minute

func (f *Fetcher) CheckExpireIssue() error {
	logrus.Info("start to check expire issue")
	opt := &github.IssueListByRepoOptions{}
	issues := []*github.Issue{}
	issues,err := f.client.GetIssues(opt)
	if err != nil {
		return err
	}
	if len(issues) == 0 {
		logrus.Info("there is no opened issue")
	}
	for _, issue := range issues {
		logrus.Info("start to check issue %d", issue.GetNumber())
		deadline := issue.GetUpdatedAt().Add(EXPIRETIME)
		current := time.Now()
		diff := current.Unix() - deadline.Unix()
		if diff > 0 {
			logrus.Info("issue %d need to be closed", issue.GetNumber())
			CloseIssue(issue)
		}
	}
	return nil
}

func CloseIssue(issue *github.Issue) error{
	jsonStr := []byte("{\"state\":\"closed\"}")
	url := "https://api.github.com/repos/hozart/pouch/issues/" + strconv.Itoa(issue.GetNumber())
	req,err  := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "token cf7d9d031891a64c1d7f4782d40fedc5237d2135")
	req.Header.Add("Content-Type", "application/json")
	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		return err
	}
	return nil
}
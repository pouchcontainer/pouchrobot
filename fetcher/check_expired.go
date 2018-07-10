package fetcher

import (
	"github.com/sirupsen/logrus"
	"github.com/google/go-github/github"
	"github.com/alibaba/pouch/pkg/utils"
	"time"
	"strconv"
	"net/http"
	"bytes"
	"github.com/pouchcontainer/pouchrobot/config"
)

const EXPIRETIME = 1 * utils.Minute
const BASEURL = "https://api.github.com/"

func (f *Fetcher) CheckExpireIssue(config config.Config) error {

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
		logrus.Info("start to check issue #", issue.GetNumber())
		deadline := issue.GetUpdatedAt().Add(EXPIRETIME)
		current := time.Now()
		diff := current.Unix() - deadline.Unix()
		if diff > 0 {
			logrus.Info("issue #"+ strconv.Itoa(issue.GetNumber()) + " need to be closed")
			CloseIssue(issue,f.client.Repo(),f.client.Owner(),config.AccessToken)
		} else {
			logrus.Info("issue #"+ strconv.Itoa(issue.GetNumber()) + "  don't need to be closed")
		}
	}
	return nil
}

func CloseIssue(issue *github.Issue,repo,owner,token string) error{
	jsonStr := []byte("{\"state\":\"closed\"}")
	url := BASEURL + "repos/"+owner+"/"+repo+"/issues/" + strconv.Itoa(issue.GetNumber())
	req,err  := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "token " + token)
	req.Header.Add("Content-Type", "application/json")
	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		logrus.Info("access " + url +" error")
		return err
	}
	return nil
}
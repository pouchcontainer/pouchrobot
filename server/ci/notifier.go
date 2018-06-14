// Copyright 2018 The Pouch Robot Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ci

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/pouchcontainer/pouchrobot/server/gh"
	"github.com/pouchcontainer/pouchrobot/server/utils"

	"github.com/google/go-github/github"
	"github.com/sirupsen/logrus"
)

// Notifier is a processor that receives notification from CI system
// and act to the messages to periodically get elements from github
type Notifier struct {
	client *gh.Client
}

// New initializes a brand new notifier
func New(client *gh.Client) *Notifier {
	return &Notifier{
		client: client,
	}
}

// TravisCIProcess gets the json string and acts to these messages from travisCI.
func (n *Notifier) TravisCIProcess(input string) error {
	input = strings.Replace(input, `\"`, `"`, -1)
	logrus.Info(input)
	var tw TravisWebhook
	if err := json.Unmarshal([]byte(input), &tw); err != nil {
		return err
	}

	prNum := tw.PullRequestNumber
	if prNum <= 0 {
		return fmt.Errorf("invalid pull request number %d unmarshalled", prNum)
	}

	logrus.Infof("TravisCI notification from PR %d received, state: %s", prNum, tw.State)

	// if the status is passed, we need to remove failure comment
	if tw.State == "passed" {
		return n.client.RmCommentsViaStr(prNum, utils.CIFailsCommentSubStr)
	}

	// if the status is failure, we need to do steps by:
	// 1. remove failure comments if there are any;
	// 2. add new failure comments to show failure state.
	if tw.State == "failed" {
		// first remove failure comments if there are any.
		n.client.RmCommentsViaStr(prNum, utils.CIFailsCommentSubStr)

		pr, err := n.client.GetSinglePR(prNum)
		if err != nil {
			return err
		}

		if pr.State != nil && *(pr.State) != "open" {
			// we only consider pr which are open
			return nil
		}
		// add new failure comments
		return n.addCIFaiureComments(pr, tw)
	}

	return nil
}

// CircleCIProcess gets the json string and acts to these messages from CircleCI.
func (n *Notifier) CircleCIProcess(input string) error {
	input = strings.Replace(input, `\"`, `"`, -1)
	logrus.Info(input)
	var cw CircleCIWebhook
	if err := json.Unmarshal([]byte(input), &cw); err != nil {
		return err
	}

	// branch in format of "branch" : "pull/1526"
	branch := cw.Branch
	prNum, err := strconv.Atoi(branch[5:])
	if err != nil {
		return fmt.Errorf("failed to get pull request number %v", err)
	}

	if prNum <= 0 {
		return fmt.Errorf("invalid pull request number %d unmarshalled", prNum)
	}

	logrus.Infof("CirleCI notification from PR %d received, state: %s", prNum, cw.Status)

	// if the status is success, we need to remove failure comment
	if cw.Status == "success" {
		return n.client.RmCommentsViaStr(prNum, utils.CIFailsCommentSubStr)
	}

	if cw.Status == "failed" {
		// first remove failure comments if there are any.
		n.client.RmCommentsViaStr(prNum, utils.CIFailsCommentSubStr)

		pr, err := n.client.GetSinglePR(prNum)
		if err != nil {
			return err
		}

		if pr.State != nil && *(pr.State) != "open" {
			// we only consider pr which are open
			return nil
		}
		// add new failure comments
		return n.addCIFaiureComments(pr, cw)
	}

	return nil
}

func (n *Notifier) addCIFaiureComments(pr *github.PullRequest, wh interface{}) error {
	// add a brand new one CI failure comments
	body := fmt.Sprintf(utils.CIFailsComment, *(pr.User.Login))
	switch wh.(type) {
	case TravisWebhook:
		tw, _ := wh.(TravisWebhook)
		detailsStr := fmt.Sprintf("build url: %s\nbuild duration: %ds\n", tw.BuildURL, tw.Duration)
		body = body + "\n" + detailsStr
	case CircleCIWebhook:
		cw, _ := wh.(CircleCIWebhook)
		detailsStr := fmt.Sprintf("build url: %s\nbuild duration: %ds\n", cw.BuildURL, cw.BuildTimeMillis)
		body = body + "\n" + detailsStr

	}

	return n.client.RmCommentsViaStrAndAttach(*(pr.Number), utils.CIFailsCommentSubStr, body)
}

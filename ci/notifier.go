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
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/pouchcontainer/pouchrobot/gh"
	"github.com/pouchcontainer/pouchrobot/utils"

	"github.com/google/go-github/github"
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

// Process gets the json string and acts to these messages from CI system, such as travisCI.
func (n *Notifier) Process(input string) error {
	input = strings.Replace(input, `\"`, `"`, -1)
	logrus.Info(input)
	var wh Webhook
	if err := json.Unmarshal([]byte(input), &wh); err != nil {
		return err
	}

	prNum := wh.PullRequestNumber
	if prNum <= 0 {
		return fmt.Errorf("invalid pull request number %d unmarshalled", prNum)
	}

	logrus.Infof("CI notification from PR %d received, state: %s", prNum, wh.State)

	// if the status is passed, we need to remove failure comment
	if wh.State == "passed" {
		return n.client.RmCommentsViaStr(prNum, utils.CIFailsCommentSubStr)
	}

	// if the status is failure, we need to do steps by:
	// 1. remove failure comments if there are any;
	// 2. add new failure comments to show failure state.
	if wh.State == "failed" {
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
		return n.addCIFailureComments(pr, wh)
	}

	return nil
}

func (n *Notifier) addCIFailureComments(pr *github.PullRequest, wh Webhook) error {
	// add a brand new one CI failure comments
	body := fmt.Sprintf(utils.CIFailsComment, *(pr.User.Login))
	detailsStr := fmt.Sprintf("build url: %s\nbuild duration: %ds\n", wh.BuildURL, wh.Duration)
	body = body + "\n" + detailsStr

	return n.client.RmCommentsViaStrAndAttach(*(pr.Number), utils.CIFailsCommentSubStr, body)
}

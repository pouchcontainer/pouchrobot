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

package pullRequestProcessor

import (
	"fmt"
	"regexp"

	"github.com/pouchcontainer/pouchrobot/server/gh"
	"github.com/pouchcontainer/pouchrobot/server/utils"

	"github.com/sirupsen/logrus"
)

var (
	// SizePrefix means the prefix of a size label
	dcoRegex = regexp.MustCompile("(?m)(Docker-DCO-1.1-)?Signed-off-by: ([^<]+) <([^<>@]+@[^<>]+)>( \\(github: ([a-zA-Z0-9][a-zA-Z0-9-]+)\\))?")
)

// PullRequestProcessor is
type PullRequestProcessor struct {
	Client *gh.Client
}

// Process processes pull request events
func (prp *PullRequestProcessor) Process(data []byte) error {
	// process details
	actionType, err := utils.ExtractActionType(data)
	if err != nil {
		return err
	}

	logrus.Infof("received event type [pull request], action type [%s]", actionType)

	pr, err := utils.ExactPR(data)
	if err != nil {
		return err
	}
	logrus.Debugf("pull request: %v", pr)

	switch actionType {
	case "opened":
		if err := prp.ActToPROpened(&pr); err != nil {
			return err
		}
	case "labeled":
		if err := prp.ActToPRLabeled(&pr); err != nil {
			return err
		}
	case "review_requested":
	case "synchronize":
		if err := prp.ActToPRSynchronized(&pr); err != nil {
			return err
		}
	case "edited":
		if err := prp.ActToPREdited(&pr); err != nil {
			return err
		}
	case "pull_request_review":
	default:
		return fmt.Errorf("unknown action type %s in pull request: ", actionType)
	}
	return nil
}

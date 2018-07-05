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

package reporter

import (
	"time"

	"github.com/pouchcontainer/pouchrobot/server/gh"
	"github.com/sirupsen/logrus"
)

// Reporter is a reporter to report weekly update on Github Repo in issues.
type Reporter struct {
	client *gh.Client
}

// New initializes a brand new reporter.
func New(client *gh.Client) *Reporter {
	return &Reporter{
		client: client,
	}
}

// Run starts to work on reporting things for repo.
func (r *Reporter) Run() {
	logrus.Infof("start to run reporter")
	// Wait time goes to Thursday.
	for {
		if time.Now().Weekday().String() == "Thursday" {
			if hour, _, _ := time.Now().Clock(); hour == 8 {
				break
			}
		}
		time.Sleep(1 * time.Hour)
	}

	for {
		// only Monday, code will enter this for loop block.
		go r.weeklyReport()

		// report one issue every week.
		time.Sleep(7 * 24 * time.Hour)
	}
}

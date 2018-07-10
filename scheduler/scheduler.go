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

package scheduler

import (
	"time"

	"github.com/pouchcontainer/pouchrobot/config"
	"github.com/pouchcontainer/pouchrobot/gh"
	"github.com/pouchcontainer/pouchrobot/scheduler/issueTask"
	"github.com/sirupsen/logrus"
)

// Scheduler for .
type Scheduler struct {
	IssueTask *issueTask.IssueTask
}

// New initializes a brand new Scheduler.
func New(client *gh.Client, config config.Config) *Scheduler {
	return &Scheduler{
		IssueTask: &issueTask.IssueTask{
			Client:                client,
			MaxDayOfNoActionIssue: config.MaxDayOfNoActionIssue,
			MaxRetryOfScheduler:   config.MaxRetryOfScheduler,
		},
	}
}

// Run starts to work on scheduling things for repo.
func (s *Scheduler) Run() {
	logrus.Infof("start to run Scheduler")
	// Wait time goes to 6 o'clock at morning every.
	for {
		if hour, _, _ := time.Now().Clock(); hour == 6 {
			break
		}
		time.Sleep(1 * time.Hour)
	}

	for {
		// only 6 o'clock of a day, code will enter this for loop block.
		go s.IssueTask.CloseOutOfDateIssues(s.IssueTask.MaxRetryOfScheduler)

		// run this everyday.
		time.Sleep(24 * time.Hour)
	}
}

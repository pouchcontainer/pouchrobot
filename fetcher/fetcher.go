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

package fetcher

import (
	"time"

	"github.com/pouchcontainer/pouchrobot/gh"

	"github.com/sirupsen/logrus"
)

// FetchInterval refers the interval of fetch action
const FetchInterval = 3 * time.Minute

// DefaultCommitGap refers default gap commit number.
// If the PR is hehind the master branch more DefaultCommitGap commits,
// robot would tell submitter to rebase.
const DefaultCommitGap = 20

// Fetcher is a worker to periodically get elements from github.
type Fetcher struct {
	client     *gh.Client
	gapCommits int
}

// New initializes a brand new fetch.
func New(client *gh.Client, CommitsGap int) *Fetcher {
	fetcher := &Fetcher{
		client:     client,
		gapCommits: CommitsGap,
	}
	if CommitsGap == 0 {
		fetcher.gapCommits = DefaultCommitGap
	}
	return fetcher
}

// Run starts periodical work
func (f *Fetcher) Run() {
	logrus.Info("start to run fetcher")

	for {
		f.CheckPRsConflict()
		f.CheckPRsGap()
		time.Sleep(FetchInterval)
	}
}

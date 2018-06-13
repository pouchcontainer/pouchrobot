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

	"github.com/sirupsen/logrus"

	"github.com/pouchcontainer/pouchrobot/server/gh"
)

// FETCHINTERVAL refers the interval of fetch action
const FETCHINTERVAL = 3 * time.Minute

// Fetcher is a worker to periodically get elements from github.
type Fetcher struct {
	client *gh.Client
}

// New initializes a brand new fetch.
func New(client *gh.Client) *Fetcher {
	return &Fetcher{
		client: client,
	}
}

// Run starts periodical work
func (f *Fetcher) Run() {
	logrus.Info("start to run fetcher")

	go func() {
		f.CheckPRsConflict()
		time.Sleep(FETCHINTERVAL)
	}()

	go func() {
		f.CheckPRsGap()
		time.Sleep(FETCHINTERVAL)
	}()
}

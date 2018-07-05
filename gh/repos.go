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

package gh

import (
	"context"

	"github.com/google/go-github/github"
	"github.com/sirupsen/logrus"
)

// GetRepository gets a repository.
func (c *Client) GetRepository() (*github.Repository, error) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	repo, _, err := c.Repositories.Get(context.Background(), c.owner, c.repo)
	if err != nil {
		logrus.Errorf("failed to get repository c.repo %s: %v", c.repo, err)
		return nil, err
	}
	logrus.Debugf("succeed in getting repository %s", c.repo)
	return repo, nil
}

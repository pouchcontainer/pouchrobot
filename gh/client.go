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
	"net/http"
	"sync"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// Client refers to a client which wishes to connect to specific repository of a user.
type Client struct {
	sync.Mutex
	*github.Client
	owner string
	repo  string
	timeUnit string
	time int
}

// NewClient constructs a new instance of Client.
func NewClient(owner, repo, timeUnit string, time int, token string) *Client {
	var tc *http.Client
	if token != "" {
		ctx := context.Background()
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{
				AccessToken: token,
			},
		)
		tc = oauth2.NewClient(ctx, ts)
	}

	return &Client{
		Client: github.NewClient(tc),
		owner:  owner,
		repo:   repo,
		timeUnit: timeUnit,
		time: time,
	}
}

// Owner returns owner of client.
func (c *Client) Owner() string {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	return c.owner
}

// Repo returns repo name of client.
func (c *Client) Repo() string {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	return c.repo
}

// Repo returns repo name of client.
func (c *Client) TimeUnit() string {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	return c.timeUnit
}

// Repo returns repo name of client.
func (c *Client) Time() int {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	return c.time
}

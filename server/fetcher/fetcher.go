package fetcher

import (
	"time"

	"github.com/allencloud/automan/server/gh"
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
	for {
		f.CheckPRsConflict()
		time.Sleep(FETCHINTERVAL)
	}
}

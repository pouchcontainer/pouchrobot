package reporter

import (
	"github.com/allencloud/automan/server/gh"
	"github.com/sirupsen/logrus"
)

// Reporter is a reporter to report weekly update on Github Repo in issues.
type Reporter struct {
	client *gh.Client
}

// New initializes a brand new reporter
func New(client *gh.Client) *Reporter {
	return &Reporter{
		client: client,
	}
}

// Run starts to work on reporting things for repo.
func (r *Reporter) Run() {
	logrus.Infof("start to run reporter")
	r.weeklyReport()
}

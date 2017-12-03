package reporter

import (
	"time"

	"github.com/allencloud/automan/server/gh"
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
	// Wait time goes to Monday.
	for {
		if time.Now().Weekday().String() == "Monday" {
			break
		}
		time.Sleep(4 * time.Hour)
	}

	for {
		if time.Now().Weekday().String() == "Monday" {
			go r.weeklyReport()
		}
		// report one issue every week.
		time.Sleep(7 * 24 * time.Hour)
	}
}

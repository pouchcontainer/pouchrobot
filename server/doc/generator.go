package doc

import (
	"time"

	"github.com/allencloud/automan/server/gh"
	"github.com/sirupsen/logrus"
)

// Generator is a processor that periodically auto generated cli and api docs for github repo.
type Generator struct {
	client *gh.Client
}

// New initializes a brand new doc generator
func New(client *gh.Client) *Generator {
	return &Generator{
		client: client,
	}
}

// Run starts periodical work
// currently generator generates doc every day.
func (g *Generator) Run() error {
	logrus.Infof("start to run doc generator")
	// Wait time goes to Monday.
	for {
		if time.Now().Weekday().String() == "Monday" {
			break
		}
		time.Sleep(4 * time.Hour)
	}

	for {
		if time.Now().Weekday().String() == "Monday" {
			go g.generateDoc()
		}
		// report one issue every two days.
		time.Sleep(2 * 24 * time.Hour)
	}
}

// generateDoc starts to generate all docs.
func (g *Generator) generateDoc() {
	// auto generate API docs
	if err := g.generateAPIDoc(); err != nil {
		logrus.Errorf("failed to generate API doc: %v", err)
	}

	// auto generate Cli docs
	if err := g.generateCliDoc(); err != nil {
		logrus.Errorf("failed to generate cli doc: %v", err)
	}
}

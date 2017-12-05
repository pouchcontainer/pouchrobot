package doc

import (
	"github.com/allencloud/automan/server/gh"
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
func (g *Generator) Run() error {
	return nil
}

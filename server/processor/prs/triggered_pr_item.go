package processor

import "github.com/allencloud/automan/server/gh"

type triggeredPRProcessor struct {
	client gh.Client
}

func (fPP *triggeredPRProcessor) Process() error {
	// process details
	return nil
}

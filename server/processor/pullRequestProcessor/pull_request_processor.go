package pullRequestProcessor

import (
	"fmt"
	"regexp"

	"github.com/allencloud/automan/server/gh"
	"github.com/allencloud/automan/server/utils"

	"github.com/sirupsen/logrus"
)

var (
	// SizePrefix means the prefix of a size label
	dcoRegex = regexp.MustCompile("(?m)(Docker-DCO-1.1-)?Signed-off-by: ([^<]+) <([^<>@]+@[^<>]+)>( \\(github: ([a-zA-Z0-9][a-zA-Z0-9-]+)\\))?")
)

// PullRequestProcessor is
type PullRequestProcessor struct {
	Client *gh.Client
}

// Process processes pull request events
func (prp *PullRequestProcessor) Process(data []byte) error {
	// process details
	actionType, err := utils.ExtractActionType(data)
	if err != nil {
		return err
	}

	logrus.Infof("received event type [pull request], action type [%s]", actionType)

	pr, err := utils.ExactPR(data)
	if err != nil {
		return err
	}
	logrus.Debugf("pull request: %v", pr)

	switch actionType {
	case "opened":
		if err := prp.ActToPROpened(&pr); err != nil {
			return err
		}
	case "labeled":
		if err := prp.ActToPRLabeled(&pr); err != nil {
			return err
		}
	case "review_requested":
	case "synchronize":
		if err := prp.ActToPRSynchronized(&pr); err != nil {
			return err
		}
	case "edited":
		if err := prp.ActToPREdited(&pr); err != nil {
			return err
		}
	case "pull_request_review":
	default:
		return fmt.Errorf("unknown action type %s in pull request: ", actionType)
	}
	return nil
}

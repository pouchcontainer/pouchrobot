package ci

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/allencloud/automan/server/gh"
	"github.com/allencloud/automan/server/utils"

	"github.com/google/go-github/github"
)

// Notifier is a processor that receives notification from CI system
// and act to the messages to periodically get elements from github
type Notifier struct {
	client *gh.Client
}

// New initializes a brand new notifier
func New(client *gh.Client) *Notifier {
	return &Notifier{
		client: client,
	}
}

// Process gets the json string and acts to these messages from CI system, such as travisCI.
func (n *Notifier) Process(input string) error {
	logrus.Info(input)
	var wh Webhook
	if err := json.Unmarshal([]byte(input), &wh); err != nil {
		return err
	}

	prNum := wh.PullRequestNumber
	if prNum <= 0 {
		return fmt.Errorf("invalid pull request number %d unmarshalled", prNum)
	}

	pr, err := n.client.GetSinglePR(prNum)
	if err != nil {
		return err
	}

	return n.addCIFaiureComments(pr, wh)
}

func (n *Notifier) addCIFaiureComments(pr *github.PullRequest, wh Webhook) error {
	// Remove all the existing CI failure comments
	comments, err := n.client.ListComments(*(pr.Number))
	if err != nil {
		return err
	}

	// remove all the existing CI failure comments
	for _, comment := range comments {
		if comment.Body == nil {
			continue
		}

		if !strings.Contains(*(comment.Body), utils.CIFailsCommentSubStr) {
			continue
		}

		n.client.RemoveComment(*(comment.ID))
	}

	// add a brand new one CI failure comments
	body := fmt.Sprintf(utils.CIFailsComment, *(pr.User.Login))
	detailsStr := fmt.Sprintf("build url: %s\nbuild duration: %ds\n", wh.BuildURL, wh.Duration)
	body = body + "\n" + detailsStr
	newComment := &github.IssueComment{
		Body: &body,
	}

	return n.client.AddCommentToPR(*(pr.Number), newComment)
}

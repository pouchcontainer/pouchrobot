package pullRequestProcessor

import (
	"github.com/google/go-github/github"
)

// ActToPRLabeled acts the event of pull request labeled.
func (prp *PullRequestProcessor) ActToPRLabeled(pr *github.PullRequest) error {
	return nil
}

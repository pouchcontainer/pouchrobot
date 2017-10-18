package fetcher

import (
	"context"
	"time"

	"github.com/google/go-github/github"
	"github.com/sirupsen/logrus"
)

// FETCHINTERVAL refers the interval of fetch action
const FETCHINTERVAL = 24 * time.Hour

// StartFetcherProcess starts to interact with GitHub every 24 hours to fetch items.
func StartFetcherProcess() {
	logrus.Info("Starting fetcher process to get github items periodically")

	//ticker := time.NewTicker(FETCHINTERVAL)

	//for {
	//	select {
	//	case <-ticker.C:
	// is there any possibility that items are so many that it consumes so much memory?
	// TODO split all below into a single goroutine to avoid that all these things cannot
	// be done in a single period.
	issues, _, err := fetchGitHubItems()
	if err != nil {
		logrus.Errorf("failed to fetch GitHub items: %v", err)
		// TODO: add a retry here
		//continue
	}

	logrus.Infof("get %d issues from the repo", len(issues))
	for _, issue := range issues {
		processSingleItem(issue)
	}

	//for _, pr := range pullRequests {
	//pr.processSingleItem()
	//	}
	//}
	//}
}

func fetchGitHubItems() ([]*github.Issue, []*github.PullRequest, error) {
	client := github.NewClient(nil)

	ctx := context.Background()
	// fetch opened issues
	issues, _, err := client.Issues.ListByRepo(ctx, "allencloud", "daoker", &github.IssueListByRepoOptions{})
	if err != nil {
		logrus.Errorf("failed to fetch issue list: %v", err)
	}

	// fetch opened pull request
	owner := "allencloud"
	repo := "daoker"
	pullRequests, _, err := client.PullRequests.List(ctx, owner, repo, &github.PullRequestListOptions{})
	if err != nil {
		logrus.Errorf("failed to fetch pull request list: %v", err)
	}

	return issues, pullRequests, nil
}

func processSingleItem(issue *github.Issue) error {
	client := github.NewClient(nil)
	ctx := context.Background()

	labels, _, err := client.Issues.ListLabels(ctx, "allencloud", "daoker", nil)
	if err != nil {
		logrus.Errorf("failed to list issue labels: %v", err)
		return err
	}
	for _, label := range labels {
		logrus.Info("%v", label)
	}
	return nil
}

package reporter

import (
	"fmt"
	"time"

	"github.com/google/go-github/github"
	"github.com/sirupsen/logrus"
)

// CalculateReviews calculates reviews from user since the last week.
func (r *Reporter) CalculateReviews(wr *WeekReport) {
	var prNums []int

	logrus.Info("start to calculate pull request reviews")
	logrus.Info("start to fetch all merged and closed pull requests")
	// fetch closed and merged pull requests since last week and all opening ones
	// first, fetch merged pull requests.
	// FIXME: we need to add closed pull requests reviews.
	query := fmt.Sprintf("is:merged type:pr repo:%s/%s merged:>=%s", r.client.Owner(), r.client.Repo(), wr.StartDate)
	issueSearchResult, err := r.client.SearchIssues(query, nil, true)
	if err != nil {
		logrus.Errorf("failed to get all merged pull requests via issue filtering: %v", err)
	}

	// SearchIssues returns a list of issue, and we can treat them as pull request as well.
	for _, pr := range issueSearchResult.Issues {
		prNums = append(prNums, *pr.Number)
	}

	// second, fetch all opening pull request
	logrus.Info("start to fetch all opening pull requests")
	if openingPRs, err := r.client.GetPullRequests(&github.PullRequestListOptions{}); err != nil {
		logrus.Errorf("failed to list all opening pull request: %v", err)
	} else {
		for _, openingPR := range openingPRs {
			prNums = append(prNums, *openingPR.Number)
		}
	}

	// add for log debuging
	logrus.Infof("these pull request %v should be taken into consideration", prNums)

	allPRReviews := []*github.PullRequestReview{}
	// get all reviews on the each of above pull request
	logrus.Info("start to get reviews from all pull requests")
	for _, prNum := range prNums {
		prReviews, err := r.client.ListPRReviews(prNum)
		if err != nil {
			logrus.Errorf("failed to get reviews from pul request %d: %v", prNum, err)
		}
		allPRReviews = append(allPRReviews, prReviews...)
	}

	prReviewsByUser := map[string]int{}
	// classify all the reviews into a map whose key is the user's ID
	for _, prReview := range allPRReviews {
		// if the review is not committed in the last week, just ignore.
		summittedAt := prReview.GetSubmittedAt()
		// FIXME: change this hard coding
		lastWeekStartAt := time.Now().Add(-7 * 24 * time.Hour)
		if summittedAt.Before(lastWeekStartAt) {
			continue
		}

		user := prReview.User
		githubID := *user.Login
		if _, exist := prReviewsByUser[githubID]; exist {
			prReviewsByUser[githubID]++
		} else {
			prReviewsByUser[githubID] = 1
		}
	}

	// set PRReviewsByUser in WeekReport
	wr.PRReviewsByUser = prReviewsByUser

	logrus.Info("succeed in calculating pull request reviews, see: %v", prReviewsByUser)
}

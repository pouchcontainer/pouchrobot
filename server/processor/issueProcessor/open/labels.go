package open

import (
	"strings"

	"github.com/pouchcontainer/pouchrobot/server/utils"

	"github.com/google/go-github/github"
	"github.com/sirupsen/logrus"
)

// ParseToGenerateLabels parses issue title and issue body to generate a slice
// with no element duplicated.
func ParseToGenerateLabels(issue *github.Issue) []string {
	var labels []string
	labels = append(labels, ParseTitleToGenerateLabels(issue)...)
	labels = append(labels, ParseBodyToGenerateLabels(issue)...)

	return utils.UniqueElementSlice(labels)
}

// ParseTitleToGenerateLabels parses issue title to generate a slice.
func ParseTitleToGenerateLabels(issue *github.Issue) []string {
	if issue.Title == nil {
		logrus.Errorf("issue %d has no title when ParseTitleToGenerateLabels", *(issue.Number))
		return nil
	}
	var labels []string
	title := issue.Title
	for label, matchedSlice := range utils.TitleMatches {
		for _, pattern := range matchedSlice {
			lowerCaseTitle := strings.ToLower(*title)
			if strings.Contains(lowerCaseTitle, pattern) {
				labels = append(labels, label)
				break
			}
		}
	}
	return labels
}

// ParseBodyToGenerateLabels parses issue title to generate a slice.
func ParseBodyToGenerateLabels(issue *github.Issue) []string {
	if issue.Body == nil {
		logrus.Errorf("issue %d has no body when ParseBodyToGenerateLabels", *(issue.Number))
		return nil
	}
	var labels []string
	content := issue.Body
	for label, matchedSlice := range utils.BodyMatches {
		for _, pattern := range matchedSlice {
			lowerCaseBody := strings.ToLower(*content)
			if strings.Contains(lowerCaseBody, pattern) {
				labels = append(labels, label)
				break
			}
		}
	}
	return labels
}

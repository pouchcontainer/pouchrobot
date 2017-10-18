package open

import (
	"strings"

	"github.com/google/go-github/github"
)

// ParseToGenerateLabels parses
func ParseToGenerateLabels(issue *github.Issue) []string {
	var labels []string
	labels = append(labels, ParseTitleToGenerateLabels(*issue)...)
	labels = append(labels, ParseBodyToGenerateLabels(*issue)...)
	return labels
}

// ParseTitleToGenerateLabels parses
func ParseTitleToGenerateLabels(issue github.Issue) []string {
	var labels []string
	title := issue.Title
	for label, matchedSlice := range TitleMatches {
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

// ParseBodyToGenerateLabels parses
func ParseBodyToGenerateLabels(issue github.Issue) []string {
	var labels []string
	content := issue.Body
	for label, matchedSlice := range BodyMatches {
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

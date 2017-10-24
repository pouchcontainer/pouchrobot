package open

import (
	"strings"

	putils "github.com/allencloud/automan/server/processor/utils"

	"github.com/google/go-github/github"
)

// ParseToGenerateLabels parses
func ParseToGenerateLabels(issue *github.Issue) []string {
	var labels []string
	labels = append(labels, ParseTitleToGenerateLabels(*issue)...)
	labels = append(labels, ParseBodyToGenerateLabels(*issue)...)

	dataMap := make(map[string]struct{}, len(labels))
	for _, value := range labels {
		if _, exist := dataMap[value]; !exist {
			dataMap[value] = struct{}{}
		}
	}
	labels = []string{}
	for key := range dataMap {
		labels = append(labels, key)
	}
	return labels
}

// ParseTitleToGenerateLabels parses
func ParseTitleToGenerateLabels(issue github.Issue) []string {
	var labels []string
	title := issue.Title
	if title == nil {
		return labels
	}
	for label, matchedSlice := range putils.TitleMatches {
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
	if content == nil {
		return labels
	}
	for label, matchedSlice := range putils.BodyMatches {
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

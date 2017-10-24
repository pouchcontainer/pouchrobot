package open

import (
	"strings"

	putils "github.com/allencloud/automan/server/processor/utils"
	"github.com/google/go-github/github"
)

var (
	// XS is
	XS = 10
	// S is
	S = 20
	// M is
	M = 40
	// L is
	L = 80
	// XL is
	XL = 160
)

// ParseToGeneratePRLabels parses
func ParseToGeneratePRLabels(pr *github.PullRequest) []string {
	var labels []string
	labels = append(labels, ParseToGetPRSize(pr))
	labels = append(labels, ParseTitleToGenerateLabels(pr)...)
	return labels
}

// ParseToGetPRSize parses the pr additions and deletions
func ParseToGetPRSize(pr *github.PullRequest) string {
	if *(pr.Additions)+*(pr.Deletions) <= XS {
		return "size/XS"
	}
	if *(pr.Additions)+*(pr.Deletions) <= S {
		return "size/S"
	}
	if *(pr.Additions)+*(pr.Deletions) <= M {
		return "size/M"
	}
	if *(pr.Additions)+*(pr.Deletions) <= L {
		return "size/L"
	}
	if *(pr.Additions)+*(pr.Deletions) <= XL {
		return "size/XL"
	}
	return "size/XXL"
}

// ParseTitleToGenerateLabels parses
func ParseTitleToGenerateLabels(pr *github.PullRequest) []string {
	var labels []string
	title := pr.Title
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

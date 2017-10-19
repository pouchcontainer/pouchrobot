package issueProcessor

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/allencloud/automan/server/gh"
	"github.com/allencloud/automan/server/processor/issueProcessor/open"
	"github.com/allencloud/automan/server/util"
	"github.com/google/go-github/github"
)

// IssueOpenedActionFunc defines
type IssueOpenedActionFunc func(issue *github.Issue) []string

// IssueOpenedActionFuncs defines
var IssueOpenedActionFuncs []IssueOpenedActionFunc

// Register registers IssueOpenedActionFunc
func Register(a IssueOpenedActionFunc) {
	IssueOpenedActionFuncs = append(IssueOpenedActionFuncs, a)
}

func init() {
	funcs := []IssueOpenedActionFunc{
		open.ParseToGenerateLabels,
	}

	for _, processFunc := range funcs {
		Register(processFunc)
	}
}

// TriggeredIssueProcessor is
type TriggeredIssueProcessor struct {
	Client *gh.Client
}

// Process processes
func (fIP *TriggeredIssueProcessor) Process(data []byte) error {
	// process details
	actionType, err := ExtractActionType(data)
	if err != nil {
		return err
	}

	logrus.Infof("received event type [issues], action type [%s]", actionType)

	issue, err := ExactIssue(data)
	if err != nil {
		return err
	}

	switch actionType {
	case "opened":
		if err := fIP.ActToIssueOpen(&issue); err != nil {
			return err
		}
	case "reopened":
	case "edited":
		if err := fIP.ActToIssueOpen(&issue); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown action type %s in issue: ", actionType)
	}

	return nil
}

// ActToIssueOpen acts to opened issue
// This function covers the following part:
// generate labels;
// attach comments;
// assign issue to specific user;
func (fIP *TriggeredIssueProcessor) ActToIssueOpen(issue *github.Issue) error {
	// generate labels
	labels := open.ParseToGenerateLabels(issue)
	if err := fIP.Client.AddLabelsToIssue(context.Background(), *(issue.Number), labels); err != nil {
		logrus.Errorf("failed to add labels %v to issue %d: %v", labels, *(issue.Number), err)
		return err
	}
	logrus.Infof("succeed in attaching labels %v to issue %d", labels, *(issue.Number))

	// attach comment
	newComment := &github.IssueComment{}

	// check if the title is too short or the body empty.
	if len(*(issue.Title)) < 20 {
		body := fmt.Sprintf(`
			Thanks for your contribution. 🍻 @%s 
			\nWhile we thought issue title could be more specific.
			\nPlease edit issue title intead of opening a new one.
			\nMore details, please refer to https://github.com/alibaba/pouch/blob/master/CONTRIBUTING.md
			`, *(issue.User.Login))
		newComment.Body = &body
		if err := fIP.Client.AddCommentToIssue(context.Background(), *(issue.Number), newComment); err != nil {
			logrus.Errorf("failed to add TOO SHORT TITLE comment to issue %d", *(issue.Number))
			return err
		}
		logrus.Infof("secceed in attaching TITLE TOO SHORT comment for issue %d", *(issue.Number))
		return nil
	}

	if issue.Body == nil || *(issue.Body) == "" || len(*(issue.Body)) < 50 {
		body := fmt.Sprintf(`
			Thanks for your contribution. 🍻 @%s 
			\nWhile we thought issue desciprtion should not be empty or too short.
			\nPlease edit this issue description intead of opening a new one.
			\nMore details, please refer to https://github.com/alibaba/pouch/blob/master/CONTRIBUTING.md
			`, *(issue.User.Login))
		newComment.Body = &body
		if err := fIP.Client.AddCommentToIssue(context.Background(), *(issue.Number), newComment); err != nil {
			logrus.Errorf("failed to add EMPTY OR TOO SHORT ISSUE BODY comment to issue %d", *(issue.Number))
			return err
		}
		logrus.Infof("secceed in attaching TITLE TOO SHORT comment for issue %d", *(issue.Number))
		return nil
	}

	// check if this is a P0 priority issue, if that mention maintainers.
	if util.SliceContainsElement(labels, "priority/P0") {
		body := fmt.Sprintf(`
			😱 \nThis is a **priority/P0** issue reported by @%s.
			\nSeems to be severe enough. 
			\nping @allencloud , PTAL. 
			`, *(issue.User.Login))
		newComment.Body = &body
		if err := fIP.Client.AddCommentToIssue(context.Background(), *(issue.Number), newComment); err != nil {
			logrus.Errorf("failed to add P0 comments to issue %d", *(issue.Number))
			return err
		}
	}
	logrus.Infof("secceed in attaching P0 comment for issue %d", *(issue.Number))

	return nil
}

// ExtractActionType extracts the action type.
func ExtractActionType(data []byte) (string, error) {
	var m struct {
		Action string `json:"action"`
	}
	if err := json.Unmarshal(data, &m); err != nil {
		return "", err
	}
	return m.Action, nil
}

// ExactIssue extracts the issue from request body.
func ExactIssue(data []byte) (github.Issue, error) {
	var m struct {
		Issue github.Issue `json:"issue"`
	}
	if err := json.Unmarshal(data, &m); err != nil {
		return github.Issue{}, err
	}
	return m.Issue, nil
}

// ExactPR extracts the pull request from request body.
func ExactPR(data []byte) (github.PullRequest, error) {
	var m struct {
		PullRequest github.PullRequest `json:"pull_reques"`
	}
	if err := json.Unmarshal(data, &m); err != nil {
		return github.PullRequest{}, err
	}
	return m.PullRequest, nil
}

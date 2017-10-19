package issueProcessor

import (
	"context"
	"fmt"
	"strings"

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

	logrus.Infof("received event type [issues] or [issue_comment], action type [%s]", actionType)

	issue, err := ExactIssue(data)
	if err != nil {
		return err
	}

	switch actionType {
	case "opened":
		if err := fIP.ActToIssueOpenOrEdit(&issue); err != nil {
			return err
		}
	case "edited":
		if err := fIP.ActToIssueOpenOrEdit(&issue); err != nil {
			return err
		}
	case "created":
		comment, err := ExactIssueComment(data)
		if err != nil {
			return err
		}
		if err := fIP.ActToIssueComment(&issue, &comment); err != nil {
			return err
		}
	case "reopened":
	default:
		return fmt.Errorf("unknown action type %s in issue: ", actionType)
	}

	return nil
}

// ActToIssueOpenOrEdit acts to opened issue
// This function covers the following part:
// generate labels;
// attach comments;
// assign issue to specific user;
func (fIP *TriggeredIssueProcessor) ActToIssueOpenOrEdit(issue *github.Issue) error {
	// generate labels
	labels := open.ParseToGenerateLabels(issue)
	if len(labels) != 0 {
		// only labels generated do we attach labels to issue
		if err := fIP.Client.AddLabelsToIssue(context.Background(), *(issue.Number), labels); err != nil {
			logrus.Errorf("failed to add labels %v to issue %d: %v", labels, *(issue.Number), err)
			return err
		}
		logrus.Infof("succeed in attaching labels %v to issue %d", labels, *(issue.Number))
	}

	// attach comment
	newComment := &github.IssueComment{}

	// check if the title is too short or the body empty.
	if len(*(issue.Title)) < 20 {
		body := fmt.Sprintf(`
Thanks for your contribution. 🍻  @%s 
While we thought **ISSUE TITLE** could be more specific.
Please edit issue title intead of opening a new one.
More details, please refer to https://github.com/alibaba/pouch/blob/master/CONTRIBUTING.md`,
			*(issue.User.Login))
		newComment.Body = &body
		if err := fIP.Client.AddCommentToIssue(context.Background(), *(issue.Number), newComment); err != nil {
			logrus.Errorf("failed to add TOO SHORT TITLE comment to issue %d", *(issue.Number))
			return err
		}
		logrus.Infof("secceed in attaching TITLE TOO SHORT comment for issue %d", *(issue.Number))

		labels := []string{"status/more-info-needed"}
		fIP.Client.AddLabelsToIssue(context.Background(), *(issue.Number), labels)

		return nil
	}

	if issue.Body == nil || *(issue.Body) == "" || len(*(issue.Body)) < 50 {
		body := fmt.Sprintf(`
Thanks for your contribution. 🍻  @%s 
While we thought **ISSUE DESCRIPTION** should not be empty or too short.
Please edit this issue description intead of opening a new one.
More details, please refer to https://github.com/alibaba/pouch/blob/master/CONTRIBUTING.md`,
			*(issue.User.Login))
		newComment.Body = &body
		if err := fIP.Client.AddCommentToIssue(context.Background(), *(issue.Number), newComment); err != nil {
			logrus.Errorf("failed to add EMPTY OR TOO SHORT ISSUE BODY comment to issue %d", *(issue.Number))
			return err
		}
		logrus.Infof("secceed in attaching TITLE TOO SHORT comment for issue %d", *(issue.Number))

		labels := []string{"status/more-info-needed"}
		fIP.Client.AddLabelsToIssue(context.Background(), *(issue.Number), labels)

		return nil
	}

	// check if this is a P0 priority issue, if that mention maintainers.
	if util.SliceContainsElement(labels, "priority/P0") {
		body := fmt.Sprintf(`
😱 This is a **priority/P0** issue reported by @%s.
Seems to be severe enough. 
ping @allencloud , PTAL. 
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

// ActToIssueComment acts to issue comment.
// It covers the following parts:
// assign to user if he comments `#dibs`
func (fIP *TriggeredIssueProcessor) ActToIssueComment(issue *github.Issue, comment *github.IssueComment) error {
	if comment == nil || issue == nil {
		return nil
	}

	commentUser := *(comment.User.Login)
	commentBody := *(comment.Body)
	users := []string{commentUser}

	if strings.Contains(strings.ToLower(commentBody), "#dibs") {
		if err := fIP.Client.AssignIssueToUsers(context.Background(), *(issue.Number), users); err != nil {
			logrus.Errorf("failed to asign users %v to issue %d", users, *(issue.Number))
			return err
		}
		logrus.Infof("secceed in assigning users %v for issue %d", users, *(issue.Number))
	}

	return nil
}

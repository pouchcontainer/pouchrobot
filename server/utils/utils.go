package utils

import "fmt"

// PRConflictLabel is a label which means conflict for pull request.
var PRConflictLabel = "conflict/needs-rebase"

// PriorityP0Label is a lable which represent P0 priority
var PriorityP0Label = "priority/P0"

// SizeLabelPrefix presents the prefix of size label name.
var SizeLabelPrefix = "size/"

// IssueTitleTooShortSubStr is a sub string used to construct the comment.
var IssueTitleTooShortSubStr = `While we thought **ISSUE TITLE** could be more specific, longer than 20 chars.
Please edit issue title instead of opening a new one.
More details, please refer to https://github.com/alibaba/pouch/blob/master/CONTRIBUTING.md`

// IssueTitleTooShort is a string used to construct the comment.
var IssueTitleTooShort = fmt.Sprintf("Thanks for your contribution. 🍻 @%s \n%s",
	"%s",
	IssueTitleTooShortSubStr,
)

// IssueDescriptionTooShortSubStr is a sub string used to construct the comment.
var IssueDescriptionTooShortSubStr = `While we thought **ISSUE DESCRIPTION** could be more specific, longer than 100 chars.
Here is a template at https://github.com/alibaba/pouch/blob/master/.github/ISSUE_TEMPLATE.md 
Please edit this issue description instead of opening a new one.
More details, please refer to https://github.com/alibaba/pouch/blob/master/CONTRIBUTING.md`

// IssueDescriptionTooShort is a string used to construct the comment.
var IssueDescriptionTooShort = fmt.Sprintf("Thanks for your contribution. 🍻 @%s \n%s",
	"%s",
	IssueDescriptionTooShortSubStr,
)

// PRTitleTooShortSubStr is a sub string used to construct the comment.
var PRTitleTooShortSubStr = `While we thought **PR TITLE** could be more specific, longer than 20 chars.
Please edit this PR title instead of opening a new one.
More details, please refer to https://github.com/alibaba/pouch/blob/master/CONTRIBUTING.md`

// PRTitleTooShort is a string used to construct the comment.
var PRTitleTooShort = fmt.Sprintf("Thanks for your contribution. 🍻  @%s \n%s",
	"%s",
	PRTitleTooShortSubStr,
)

// PRDescriptionTooShortSubStr is a sub string used to construct the comment.
var PRDescriptionTooShortSubStr = `While we thought **PR Description** could be more specific, longer than 100 chars.
Please edit this PR title instead of opening a new one.
More details, please refer to https://github.com/alibaba/pouch/blob/master/CONTRIBUTING.md`

// PRDescriptionTooShort is a string used to construct the comment.
var PRDescriptionTooShort = fmt.Sprintf("Thanks for your contribution. 🍻  @%s \n%s",
	"%s",
	PRDescriptionTooShortSubStr,
)

// PRNeedsSignOff is a string used to remind user to sign off.
var PRNeedsSignOff = `
Thanks for your contribution. 🍻  @%s
Please sign off in each of your commits.
`

// IssueNeedPOCommentSubStr is a string used to attach comment on P0 issue.
var IssueNeedPOCommentSubStr = `😱 This is a **priority/P0** issue.
Seems to be severe enough. 
ping @alibaba/pouch , PTAL. 
`

// IssueNeedPOComment is a string used to attach comment on P0 issue.
var IssueNeedPOComment = fmt.Sprintf("Thanks for your report, @%s \n%s",
	"%s",
	IssueNeedPOCommentSubStr,
)

// FirstCommitComment is a comment used to thank a user's first contribution.
var FirstCommitComment = `
Thanks for your first contribution, @%s
`

// PRConflictSubStr is a substring of conflict message.
var PRConflictSubStr = "Conflict happens after merging a previous commit. Please rebase the branch against master and push it back again. Thanks a lot."

// PRConflictComment is a string used to attach comment on conflict PR.
var PRConflictComment = fmt.Sprintf("ping @%s \n%s",
	"%s",
	PRConflictSubStr,
)

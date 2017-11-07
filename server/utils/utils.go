package utils

import "fmt"

// ConflictLabel is a label which means conflict for pull request
var ConflictLabel = "conflict/needs-rebase"

// IssueTitleTooShort is a string used to construct the comment
var IssueTitleTooShort = `
Thanks for your contribution. 🍻  @%s 
While we thought **ISSUE TITLE** could be more specific, longer than 20 chars.
Please edit issue title instead of opening a new one.
More details, please refer to https://github.com/alibaba/pouch/blob/master/CONTRIBUTING.md`

// IssueDescriptionTooShort is a string used to construct the comment
var IssueDescriptionTooShort = `
Thanks for your contribution. 🍻  @%s 
While we thought **ISSUE DESCRIPTION** should not be empty or too short.
Here is a template at https://github.com/alibaba/pouch/blob/master/.github/ISSUE_TEMPLATE.md 
Please edit this issue description instead of opening a new one.
More details, please refer to https://github.com/alibaba/pouch/blob/master/CONTRIBUTING.md`

// PRTitleTooShort is a string used to construct the comment
var PRTitleTooShort = `
Thanks for your contribution. 🍻  @%s 
While we thought **PR TITLE** could be more specific, longer than 20 chars.
Please edit this PR title instead of opening a new one.
More details, please refer to https://github.com/alibaba/pouch/blob/master/CONTRIBUTING.md`

// PRDescriptionTooShort is a string used to construct the comment
var PRDescriptionTooShort = `
Thanks for your contribution. 🍻  @%s 
While we thought **PR DESCRIPTION** should not be empty or too short.
Here is a template at https://github.com/alibaba/pouch/blob/master/.github/PULL_REQUEST_TEMPLATE.md.
Please edit this PR description instead of opening a new one.
More details, please refer to https://github.com/alibaba/pouch/blob/master/CONTRIBUTING.md`

// IssueNeedPOComment is a string used to attach comment on P0 issue.
var IssueNeedPOComment = `
😱 This is a **priority/P0** issue reported by @%s.
Seems to be severe enough. 
ping @alibaba/pouch , PTAL. 
`

// FirstCommitComment is a comment used to thank a user's first contribution.
var FirstCommitComment = `
Thanks for your first contribution, @%s
`

// ConflictSubStr is a substring of conflict message.
var ConflictSubStr = "Conflict happens after merging a previous commit. Please rebase the branch against master and push it back again. Thanks a lot."

// PRConflictComment is a string used to attach comment on conflict PR.
var PRConflictComment = fmt.Sprintf("ping @%s \n%s", "%s", ConflictSubStr)

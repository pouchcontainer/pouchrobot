package utils

// IssueTitleTooShort is a string used to construct the comment
var IssueTitleTooShort = `
Thanks for your contribution. 🍻  @%s 
While we thought **ISSUE TITLE** could be more specific.
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
While we thought **PR TITLE** should not be empty or too short.
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

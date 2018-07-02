// Copyright 2018 The Pouch Robot Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package utils

import "fmt"

// CIFailureLable is a label which means ci failure.
var CIFailureLable = "CI-failure"

// PRConflictLabel is a label which means conflict for pull request.
var PRConflictLabel = "conflict/needs-rebase"

//PRGapLabel is a label which means gap for pull request.
var PRGapLabel = "gap/needs-rebase"

// PriorityP1Label is a lable which represent P1 priority which is highest.
var PriorityP1Label = "priority/P1"

// SizeLabelPrefix presents the prefix of size label name.
var SizeLabelPrefix = "size/"

// IssueTitleTooShortSubStr is a sub string used to construct the comment.
var IssueTitleTooShortSubStr = `While we thought **ISSUE TITLE** could be more specific, longer than 20 chars.
Please edit issue title instead of opening a new one.
More details, please refer to https://github.com/alibaba/pouch/blob/master/CONTRIBUTING.md`

// IssueTitleTooShort is a string used to construct the comment.
var IssueTitleTooShort = fmt.Sprintf(
	"Thanks for your contribution. 🍻 @%s \n%s",
	"%s",
	IssueTitleTooShortSubStr,
)

// IssueDescriptionTooShortSubStr is a sub string used to construct the comment.
var IssueDescriptionTooShortSubStr = `While we thought **ISSUE DESCRIPTION** could be more specific, longer than 100 chars.
Here is a template at https://github.com/alibaba/pouch/blob/master/.github/ISSUE_TEMPLATE.md
Please edit this issue description instead of opening a new one.
More details, please refer to https://github.com/alibaba/pouch/blob/master/CONTRIBUTING.md`

// IssueDescriptionTooShort is a string used to construct the comment.
var IssueDescriptionTooShort = fmt.Sprintf(
	"Thanks for your contribution. 🍻 @%s \n%s",
	"%s",
	IssueDescriptionTooShortSubStr,
)

// PRTitleTooShortSubStr is a sub string used to construct the comment.
var PRTitleTooShortSubStr = `While we thought **PR TITLE** could be more specific, longer than 20 chars.
Please edit this PR title instead of opening a new one.
More details, please refer to https://github.com/alibaba/pouch/blob/master/CONTRIBUTING.md`

// PRTitleTooShort is a string used to construct the comment.
var PRTitleTooShort = fmt.Sprintf(
	"Thanks for your contribution. 🍻  @%s \n%s",
	"%s",
	PRTitleTooShortSubStr,
)

// PRDescriptionTooShortSubStr is a sub string used to construct the comment.
var PRDescriptionTooShortSubStr = `While we thought **PR Description** could be more specific, longer than 100 chars.
Please edit this PR title instead of opening a new one.
More details, please refer to https://github.com/alibaba/pouch/blob/master/CONTRIBUTING.md`

// PRDescriptionTooShort is a string used to construct the comment.
var PRDescriptionTooShort = fmt.Sprintf(
	"Thanks for your contribution. 🍻  @%s \n%s",
	"%s",
	PRDescriptionTooShortSubStr,
)

// PRNeedsSignOffStr is a string used to remind user to sign off.
var PRNeedsSignOffStr = ` Thanks for your contribution. 🍻
Please sign off in each of your commits.`

// PRNeedsSignOff is a string used to remind user to sign off.
var PRNeedsSignOff = fmt.Sprintf(
	"@%s %s",
	"%s",
	PRNeedsSignOffStr,
)

// IssueNeedP1CommentSubStr is a string used to attach comment on P1 issue.
var IssueNeedP1CommentSubStr = `😱 This is a **priority/P1** issue which is highest.
Seems to be severe enough.
ping @alibaba/pouch , PTAL.
`

// IssueNeedP1Comment is a string used to attach comment on P1 issue.
var IssueNeedP1Comment = fmt.Sprintf(
	"Thanks for your report, @%s \n%s",
	"%s",
	IssueNeedP1CommentSubStr,
)

// FirstCommitCommentSubStr is a string which is substring of FirstCommitComment
var FirstCommitCommentSubStr = `👏  We really appreciate it.
Just remind that you have read the contribution guide: https://github.com/alibaba/pouch/blob/master/CONTRIBUTING.md
If you didn't, you should do that first. If done, welcome again and please enjoy hacking! 🍻
`

// FirstCommitComment is a comment used to thank a user's first contribution.
var FirstCommitComment = fmt.Sprintf(
	"We found this is your first time to contribute to Pouch, @%s \n%s",
	"%s",
	FirstCommitCommentSubStr,
)

// PRConflictSubStr is a substring of conflict message.
var PRConflictSubStr = `Conflict happens after merging a previous commit.
Please rebase the branch against master and push it back again. Thanks a lot.
`

// PRConflictComment is a string used to attach comment on conflict PR.
var PRConflictComment = fmt.Sprintf(
	"ping @%s \n%s",
	"%s",
	PRConflictSubStr,
)

// PRGapSubStr is a substring of gap message.
var PRGapSubStr = `We found that this PR is over 10 commits behind master.
Please rebase the branch against master and push it back again. Thanks a lot.
`

// PRGapComment is a string used to attach comment on gap PR.
var PRGapComment = fmt.Sprintf(
	"ping @%s \n%s",
	"%s",
	PRGapSubStr,
)

// CIFailsCommentSubStr is a substring of CI failing comments
var CIFailsCommentSubStr = `
CI fails according integration system.
Please refer to the CI failure Details button to corresponding test, and update your PR to pass CI.

If this is flaky test, welcome to track this with [profiling an issue](https://github.com/alibaba/pouch/issues/new).
`

// CIFailsComment is a string used to attach comment to CI failed PRs.
var CIFailsComment = fmt.Sprintf(
	"ping @%s \n%s",
	"%s",
	CIFailsCommentSubStr,
)

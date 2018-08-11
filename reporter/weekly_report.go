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

package reporter

// WeekReport contains details about elements to construct a report.
type WeekReport struct {
	// time of this weekly report's start time.
	StartDate string

	// time of this weekly report's end time.
	EndDate string

	// Watch defines currently how many github users are watching this repo.
	Watch int

	// Star defines currently how many github users are staring this repo.
	Star int

	// Fork defines currently how many github users have forked this repo.
	Fork int

	// Contributors defines the number of contributors.
	Contributors int

	// NumOfNewIssues is the issues number which are created in the last week.
	NumOfNewIssues int

	// NumOfClosedIssues is the issues number which are closed in the last week.
	NumOfClosedIssues int

	// MergedPR defines how many pull requests have beem merge between time StartDate and EndDate.
	MergedPR map[string][]*SimplePR

	// CountOfPR defines the number of merged pull request.
	CountOfPR int

	// NewContributors defines new contributors between time StartDate and EndDate.
	NewContributors []string

	// PRReviewsByUser defines that all pull request reviews submitted between time StartDate and EndDate.
	// PRReviewsByUser has a type map, the key is User, Value is the number of pull reuqest reviews of single User.
	PRReviewsByUser map[string]int
}

// StatsLastWeek collects repo data from last week.
type StatsLastWeek struct {
	Watch        int
	Star         int
	Fork         int
	Contributors int
}

// SimplePR represents
type SimplePR struct {
	Num     int
	Title   string
	HTMLURL string
}

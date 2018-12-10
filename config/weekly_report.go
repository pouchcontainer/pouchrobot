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

package config

// WeeklyReportConfig refers to auto weekly report config
type WeeklyReportConfig struct {
	// ReportDay representing which is the weekly report generation day.
	ReportDay string `json:"reportDay"`

	// ReportHour representing which is the weekly report generation time on ReportDay.
	ReportHour int `json:"reportHour"`
}

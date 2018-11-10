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

// Config refers
type Config struct {
	// Owner is the organization of open source project.
	Owner string

	// Repo is the repository name.
	Repo string

	// HTTPListen is the tcp address the robot listens on.
	HTTPListen string

	// AccessToken is identify which github user this robot plays the role.
	AccessToken string

	// Commits Gap is for fetcher to check commit gap between pr and master branch,
	// if it is larger than CommitsGap, request to rebase this.
	CommitsGap int

	// For weekly reporter

	// ReportDay representing which is the weekly report generation day.
	ReportDay string

	// ReportHour representing which is the weekly report generation time on ReportDay.
	ReportHour int

	// For doc generator

	// RootDir specifies repo's rootdir which is to generated docs.
	RootDir string

	// SwaggerPath specifies where the swagger.yml file locates.
	SwaggerPath string

	// APIDocPath specifies where to generate the doc file corresponding to swagger.yml.
	// this is a relative path to root dir.
	APIDocPath string

	// GenerationHour represents doc generation time every day.
	// Valid range is [0, 23].
	GenerationHour int
}

// NewConfig creates a brand new Config instance with default values.
func NewConfig() Config {
	return Config{}
}

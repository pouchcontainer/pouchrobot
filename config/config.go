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

// CmdConfig refers the start command line config needed
type CmdConfig struct {
	// ConfigFilePath is the path of config json file
	ConfigFilePath string

	// Debug refers to the log mode.
	Debug bool
}

// Config refers the config values for the project
type Config struct {
	// Owner is the organization of open source project.
	Owner string `json:"owner"`

	// Repo is the repository name.
	Repo string `json:"repo"`

	// HTTPListen is the tcp address the robot listens on.
	HTTPListen string `json:"httpListen"`

	// AccessToken is identify which github user this robot plays the role.
	AccessToken string `json:"accessToken"`

	// FetcherConfig is configs for fetcher module
	FetcherConfig FetcherConfig `json:"fetcher"`

	// DocGenerateConfig is configs for doc generate module
	DocGenerateConfig DocGenerateConfig `json:"docGenerator"`

	// TranslatorConfig is configs for translate module
	TranslatorConfig TranslatorConfig `json:"translator"`

	// WeeklyReportConfig is configs for weekly report module
	WeeklyReportConfig WeeklyReportConfig `json:"weeklyReport"`
}

// NewConfig creates a brand new Config instance
func NewConfig() Config {
	return Config{}
}

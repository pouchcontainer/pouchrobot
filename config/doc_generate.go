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

// DocGenerateConfig refers to doc auto generate config
type DocGenerateConfig struct {
	// RootDir specifies repo's rootdir which is to generated docs.
	RootDir string `json:"rootDir"`

	// SwaggerPath specifies where the swagger.yml file locates.
	SwaggerPath string `json:"swaggerPath"`

	// APIDocPath specifies where to generate the doc file corresponding to swagger.yml.
	// this is a relative path to root dir.
	APIDocPath string `json:"APIDocPath"`

	// GenerationHour represents doc generation time every day.
	// Valid range is [0, 23].
	GenerationHour int `json:"generationHour"`
}

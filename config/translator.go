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

// TranslatorConfig refers to auto translate config
type TranslatorConfig struct {
	// BaiduConfig is the config for baidu based translator
	BaiduConfig BaiduTranslateConfig `json:"baidu"`
}

// BaiduTranslateConfig refers to Baidu API based translate config
type BaiduTranslateConfig struct {
	// AppID is the appid for baidu translator init
	AppID string `json:"appID"`

	// Key is the appid for baidu translator init
	Key string `json:"key"`
}

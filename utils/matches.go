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

// Maintainers is a list of the maintainers,
// TODO: this part will be auto-generated according to MAINTAINERS file in repo.
var Maintainers = []string{
	"allencloud",
	"yyb196",
	"Ace-Tang",
	"skoo87",
	"sunyuan3",
	"furykerry",
	"WIZARD-CXY",
	"skyline09",
	"rudyfly",
	"houqianming",
	"Letty5411",
	"HusterWan",
	"shaloulcy",
}

// TitleMatches is a map in which key is the label, and value is a slice of string
// which can be treated as the label.
var TitleMatches = map[string][]string{
	"areas/cli": {
		"cli:",
		"cli :",
		"command",
		"command line",
		"command-line",
	},
	"areas/docs": {
		"doc:",
		"docs:",
		"doc :",
		"docs :",
		"document",
	},
	"areas/log": {
		"gelf",
		"fluentd",
		"journald",
		"log",
		"splunk",
		"syslog",
	},
	"areas/images": {
		"docker image",
		"image-spec",
		"pouch pull",
	},
	"areas/monitoring": {
		"monitoring",
		"prometheus",
		"health check",
	},
	"areas/network": {
		"cni",
		"ipvlan",
		"ipsec",
		"macvlan",
		"network",
		"overlay",
		"vlan",
		"vxlan",
	},
	"areas/orchestration": {
		"kubernetes",
		"marathon",
		"mesos",
		"swarm",
		"swarmkit",
	},
	"areas/runv": {
		"runv",
	},
	"areas/storage": {
		"csi",
		"storage",
		"volume",
	},
	"areas/test": {
		"ci",
		"test",
	},
	"areas/typo": {
		"typo",
	},
	"kind/bug": {
		"bug",
		"bugfix",
		"cannot",
		"can not",
		"can't",
		"error",
		"failure",
		"failed to ",
		"fix:",
	},
	"kind/design": {
		"design",
	},
	"kind/feature": {
		"feature",
	},
	"kind/feature-request": {
		"feature request",
		"feature-request",
		"feature_request",
	},
	"kind/panic": {
		"invalid memory address or nil pointer",
		"panic",
	},
	"kind/performance": {
		"performance",
	},
	"kind/proposal": {
		"proposal",
	},
	"kind/question": {
		"can i",
		"can you",
		"confusion",
		"does pouch",
		"how to",
		"question",
		"where to",
	},
	"kind/refactor": {
		"refactor",
	},
	"os/windows": {
		"windows",
		"windows server",
		".net",
	},
	"os/ubuntu": {
		"ubuntu",
	},
	"os/macos": {
		"macos",
		"osx",
	},
	"os/centos": {
		"centos",
	},
	"os/fedora": {
		"fedora",
	},
	"os/suse": {
		"suse",
	},
	"os/freebsd": {
		"freebsd",
	},
	"priority/P1": {
		"panic",
		"invalid memory address or nil pointer",
	},
	"DO-NOT-MERGE": {
		"do not merge",
		"do-not-merge",
		"don't merge",
	},
	"WeeklyReport": {
		"weekly report",
		"weeklyreport",
		"weekreport",
		"week report",
	},
}

// BodyMatches is a map in which key is the label, and value is a slice of string
// which can be treated as the label.
var BodyMatches = map[string][]string{
	"kind/panic": {
		"panic",
		"invalid memory address or nil pointer",
	},
}

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
	"areas/cli": []string{
		"cli:",
		"cli :",
		"command",
		"command line",
		"command-line",
	},
	"areas/docs": []string{
		"doc:",
		"docs:",
		"doc :",
		"docs :",
		"document",
	},
	"areas/log": []string{
		"gelf",
		"fluentd",
		"journald",
		"log",
		"splunk",
		"syslog",
	},
	"areas/images": []string{
		"docker image",
		"image-spec",
		"pouch pull",
	},
	"areas/monitoring": []string{
		"monitoring",
		"prometheus",
		"health check",
	},
	"areas/network": []string{
		"cni",
		"ipvlan",
		"ipsec",
		"macvlan",
		"network",
		"overlay",
		"vlan",
		"vxlan",
	},
	"areas/orchestration": []string{
		"kubernetes",
		"marathon",
		"mesos",
		"swarm",
		"swarmkit",
	},
	"areas/runv": []string{
		"runv",
	},
	"areas/storage": []string{
		"csi",
		"storage",
		"volume",
	},
	"areas/test": []string{
		"ci",
		"test",
	},
	"areas/typo": []string{
		"typo",
	},
	"kind/bug": []string{
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
	"kind/design": []string{
		"design",
	},
	"kind/feature": []string{
		"feature",
	},
	"kind/feature-request": []string{
		"feature request",
		"feature-request",
		"feature_request",
	},
	"kind/panic": []string{
		"invalid memory address or nil pointer",
		"panic",
	},
	"kind/performance": []string{
		"performance",
	},
	"kind/proposal": []string{
		"proposal",
	},
	"kind/question": []string{
		"can i",
		"can you",
		"confusion",
		"does pouch",
		"how to",
		"question",
		"where to",
	},
	"kind/refactor": []string{
		"refactor",
	},
	"os/windows": []string{
		"windows",
		"windows server",
		".net",
	},
	"os/ubuntu": []string{
		"ubuntu",
	},
	"os/macos": []string{
		"macos",
		"osx",
	},
	"os/centos": []string{
		"centos",
	},
	"os/fedora": []string{
		"fedora",
	},
	"os/suse": []string{
		"suse",
	},
	"os/freebsd": []string{
		"freebsd",
	},
	"priority/P1": []string{
		"panic",
		"invalid memory address or nil pointer",
	},
	"DO-NOT-MERGE": []string{
		"do not merge",
		"do-not-merge",
		"don't merge",
	},
	"WeeklyReport": []string{
		"weekly report",
		"weeklyreport",
		"weekreport",
		"week report",
	},
}

// BodyMatches is a map in which key is the label, and value is a slice of string
// which can be treated as the label.
var BodyMatches = map[string][]string{
	"kind/panic": []string{
		"panic",
		"invalid memory address or nil pointer",
	},
}

// SpecialIssueMatches is a map in which key is the label, and value is a slice of labels
// which will be specifically attached to those issues which contain this key.
var SpecialIssueMatches = map[string][]string{
	"WeeklyReport": []string{
		"WeeklyReport",
	},
}

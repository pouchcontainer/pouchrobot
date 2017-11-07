package utils

// Maintainers is a list of the maintainers
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
}

// TitleMatches is a map in which key is the label, and value is a slice of string
// which can be treated as the label.
var TitleMatches = map[string][]string{
	"areas/cli": []string{
		"cli:",
		"cli :",
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
	"priority/P0": []string{
		"panic",
		"invalid memory address or nil pointer",
	},
	"areas/typo": []string{
		"typo",
	},
	"DO-NOT-MERGE": []string{
		"do not merge",
		"do-not-merge",
		"don't merge",
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

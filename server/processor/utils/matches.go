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
	"areas/log": []string{
		"gelf",
		"fluentd",
		"journald",
		"splunk",
		"syslog",
	},
	"areas/monitoring": []string{
		"monitoring",
		"prometheus",
	},
	"areas/network": []string{
		"cni",
		"network",
		"overlay",
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
		"cannot",
		"can not",
		"can't",
		"error",
		"failure",
		"failed to ",
	},
	"kind/panic": []string{
		"invalid memory address or nil pointer",
		"panic",
	},
	"kind/propasal": []string{
		"proposal",
	},
	"kind/design": []string{
		"design",
	},
	"kind/performance": []string{
		"performance",
	},
	"kind/feature-request": []string{
		"feature request",
		"feature-request",
		"feature_request",
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
}

// BodyMatches is a map in which key is the label, and value is a slice of string
// which can be treated as the label.
var BodyMatches = map[string][]string{
	"kind/panic": []string{
		"panic",
		"invalid memory address or nil pointer",
	},
}

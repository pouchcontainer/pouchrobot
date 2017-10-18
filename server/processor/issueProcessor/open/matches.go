package open

// TitleMatches is a map in which key is the label, and value is a slice of string
// which can be treated as the label.
var TitleMatches = map[string][]string{
	"kind/bug": []string{
		"bug",
		"error",
		"failure",
		"failed to ",
		"cannot",
		"can not",
		"can't",
	},
	"kind/panic": []string{
		"panic",
		"invalid memory address or nil pointer",
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
		"question",
		"confusion",
		"how to",
		"can i",
		"does pouch",
		"can you",
		"where to",
	},
	"areas/network": []string{
		"network",
		"cni",
		"vxlan",
		"overlay",
	},
	"areas/storage": []string{
		"volume",
		"storage",
		"csi",
	},
	"areas/orchestration": []string{
		"kubernetes",
	},
	"areas/runv": []string{
		"runv",
	},
	"areas/test": []string{
		"ci",
		"test",
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

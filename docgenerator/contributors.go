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

package docgenerator

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

// ContributorsFilePath identifies the contributors path in root dir.
var ContributorsFilePath = "CONTRIBUTORS"

// generateContributors will generate file CONTRIBUTORS.
func (g *Generator) generateContributors() error {
	// auto generate cli docs
	// cmd := exec.Command("/bin/bash", "-c", "/go/src/github.com/alibaba/pouch/hack/generate-contributors.sh")
	// if err := cmd.Run(); err != nil {
	//	return fmt.Errorf("failed to generate contributors: %v", err)
	// }

	f, err := os.Create(filepath.Join(g.RootDir, ContributorsFilePath))
	if err != nil {
		logrus.Errorf("failed to create CONTRIBUTORS file in target local repo: %v", err)
		return err
	}

	header := GetContributorsFileHeader()

	contributorsList, err := GenContributorsList()
	if err != nil {
		logrus.Errorf("failed to generate CONTRIBUTORS list by git log: %v", err)
		return err
	}

	allContent := header + contributorsList

	if _, err := f.Write([]byte(allContent)); err != nil {
		logrus.Errorf("failed to write data to CONTRIBUTORS file: %v", err)
		return err
	}

	return nil
}

// GetContributorsFileHeader will return a common header for file CONTRIBUTOR,
// and show how to get this file.
func GetContributorsFileHeader() string {
	return `# This file lists all contributors having contributed to this project.
# For how it is generated, see command "git log --format='%aN <%aE>' | sort -uf".

`
}

// GenContributorsList generates the contributors list via git command.
func GenContributorsList() (string, error) {
	cmd := exec.Command("/bin/bash", "-c", "git log --format='%aN <%aE>' | sort -uf")

	bytes, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

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
	"fmt"
	"os/exec"
)

// generateContributors will generate file CONTRIBUTORS.
// First, use newly built binary pouch to execute `generate-conrtibutors.sh` to generate Cli doc.
// Second, git commit and push to github.
// Third, use github to create a new pull request.
func (g *Generator) generateContributors() error {
	// auto generate cli docs
	cmd := exec.Command("/bin/bash", "-c", "/go/src/github.com/alibaba/pouch/hack/generate-contributors.sh")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to generate contributors: %v", err)
	}

	return nil
}

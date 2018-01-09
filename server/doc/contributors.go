package doc

import (
	"fmt"
	"os/exec"
)

// generateContributors will generate file CONTRIBUTORS.
// First, use newly built binary pouch to execute `pouch gen-doc` to generate Cli doc.
// Second, git commit and push to github.
// Third, use github to create a new pull request.
func (g *Generator) generateContributors() error {
	// auto generate cli docs
	cmd := exec.Command("/bin/bash", "hack/generate-conrtibutors.sh")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to gen doc via cobra: %v", err)
	}

	return nil
}

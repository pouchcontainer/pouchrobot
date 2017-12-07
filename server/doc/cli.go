package doc

import (
	"fmt"
	"os/exec"
)

// generateCliDoc will generate Cli doc.
// First, use newly built binary pouch to execute `pouch gen-doc` to generate Cli doc.
// Second, git commit and push to github.
// Third, use github to create a new pull request.
func (g *Generator) generateCliDoc() error {
	// build a new pouch cli client, since all cli doc is from newly built cli.
	cmd := exec.Command("make", "client")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to make client: %v", err)
	}

	// auto generate cli docs
	cmd = exec.Command("./pouch", "gen-doc")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to gen doc via cobra: %v", err)
	}

	return nil
}

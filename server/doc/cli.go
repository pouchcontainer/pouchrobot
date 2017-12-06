package doc

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/google/go-github/github"
)

// generateCliDoc will generate Cli doc.
// First, use newly built binary pouch to execute `pouch gen-doc` to generate Cli doc.
// Second, git commit and push to github.
// Third, use github to create a new pull request.
func (g *Generator) generateCliDoc() error {
	newBranch := generateNewBranch()
	logrus.Infof("generate a new branch name %s", newBranch)

	// sync latest master branch and checkout new branch
	cmd := exec.Command("git", "checkout", "master")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to checkout master: %v", err)
	}
	cmd = exec.Command("git", "fetch", "upstream", "master")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to git fetch upstreanm master: %v", err)
	}
	cmd = exec.Command("git", "rebase", "upstream/master")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to git rebase upstreanm/master: %v", err)
	}
	cmd = exec.Command("git", "push", "-f", "origin", "master")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to git push -f origin master: %v", err)
	}
	// create a new branch
	cmd = exec.Command("git", "checkout", "-b", newBranch)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to git checkout -b %s: %v", newBranch, err)
	}

	cmd = exec.Command("make", "client")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to make client: %v", err)
	}

	// auto generate cli docs
	cmd = exec.Command("./pouch", "gen-doc")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to gen doc via cobra: %v", err)
	}

	// commit and push branch
	if err := gitCommitAndPush(newBranch); err != nil {
		return err
	}

	// start to submit pull request
	if err := g.sumbitPR(newBranch); err != nil {
		return err
	}
	return nil
}

func gitCommitAndPush(newBranch string) error {
	// git add all updated files.
	cmd := exec.Command("git", "add", ",")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to git add .: %v", err)
	}

	// check whether nothing changed.
	out, err := exec.Command("git", "status").Output()
	if err != nil {
		return err
	}

	// if nothing changes, return nil to quit git procedure.
	if strings.Contains(string(out), "nothing to commit") {
		logrus.Infof("no cli doc changes happened, quit git procedure")
		return nil
	}

	// git commit all the staged files.
	cmd = exec.Command("git", "commit", "-s", "-m", "docs: auto generate pouch cli docs via code")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to git commit -s -m : %v", err)
	}

	// git push forcely to origin repo.
	cmd = exec.Command("git", "push", "-f", "origin", newBranch)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to git push -f origin %s: %v", newBranch, err)
	}

	// git branch -D to delete branch to free resources.
	cmd = exec.Command("git", "branch", "-D", newBranch)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to git push branch -D %s: %v", newBranch, err)
	}

	return nil
}

func (g *Generator) sumbitPR(branch string) error {
	title := "docs: auto generate pouch cli docs via code"
	head := fmt.Sprintf("pouchrobot:%s", branch)
	base := "master"
	body := `**1.Describe what this PR did**
	
	**2.Does this pull request fix one issue?** 
	None
	
	**3.Describe how you did it**
	None
	
	**4.Describe how to verify it**
	None
	
	**5.Special notes for reviews**
	None`

	newPR := &github.NewPullRequest{
		Title: &title,
		Head:  &head,
		Base:  &base,
		Body:  &body,
	}
	if _, err := g.client.CreatePR(newPR); err != nil {
		return err
	}
	return nil
}

func generateNewBranch() string {
	timeStr := time.Now().String()
	dateStrSlice := strings.SplitN(timeStr, " ", 2)
	dateStr := dateStrSlice[0]

	return fmt.Sprintf("cli-doc-%s", dateStr)
}

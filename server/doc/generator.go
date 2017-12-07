package doc

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/allencloud/automan/server/gh"
	"github.com/google/go-github/github"
	"github.com/sirupsen/logrus"
)

// ErrNothingChanged is used when git commit has nothing to commit.
var ErrNothingChanged = fmt.Errorf("nothing to commit")

// Generator is a processor that periodically auto generated cli and api docs for github repo.
type Generator struct {
	client *gh.Client
}

// New initializes a brand new doc generator
func New(client *gh.Client) *Generator {
	return &Generator{
		client: client,
	}
}

// Run starts periodical work
// currently generator generates doc every day.
func (g *Generator) Run() error {
	logrus.Infof("start to run doc generator")
	/*for {
		// Break the loop if the time passes one clock
		// Since time zone container maybe has a delta 8 hours from Beijing Time,
		// It is about 9 o'clock in Beijing Time.
		if hour, _, _ := time.Now().Clock(); hour == 1 {
			break
		}
		time.Sleep(30 * time.Minute)
	}*/

	for {
		go g.generateDoc()
		// generate cli and api docs every day.
		time.Sleep(24 * time.Hour)
	}
}

// generateDoc starts to generate all docs.
func (g *Generator) generateDoc() error {
	newBranch := generateNewBranch()
	logrus.Infof("generate a new branch name %s", newBranch)

	// do prepare thing before cli and api doc generation.
	if err := prepareGitEnv(newBranch); err != nil {
		logrus.Errorf("failed to prepare git environment: %v", err)
		return err
	}

	// auto generate API docs on local filesystem.
	if err := g.generateAPIDoc(); err != nil {
		logrus.Errorf("failed to generate API doc: %v", err)
	}

	// auto generate Cli docs on local filesystem.
	if err := g.generateCliDoc(); err != nil {
		logrus.Errorf("failed to generate cli doc: %v", err)
	}

	// commit and push branch
	if err := gitCommitAndPush(newBranch); err != nil {
		if err == ErrNothingChanged {
			// if nothing changed, no need to submit pull request.
			return nil
		}
		return err
	}

	// start to submit pull request
	if err := g.sumbitPR(newBranch); err != nil {
		return err
	}
	return nil
}

func prepareGitEnv(newBranch string) error {
	// sync latest master branch and checkout new branch

	// checkout local master branch
	cmd := exec.Command("git", "checkout", "master")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to checkout master: %v", err)
	}

	// fetch upstream master to local
	cmd = exec.Command("git", "fetch", "upstream", "master")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to git fetch upstreanm master: %v", err)
	}

	// rebase local master on origin/master
	cmd = exec.Command("git", "rebase", "upstream/master")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to git rebase upstreanm/master: %v", err)
	}

	// push local master to origin/master
	cmd = exec.Command("git", "push", "-f", "origin", "master")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to git push -f origin master: %v", err)
	}

	// create a new branch named by input newBranch
	// the following doc generation are all on this new branch
	cmd = exec.Command("git", "checkout", "-b", newBranch)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to git checkout -b %s: %v", newBranch, err)
	}
	return nil
}

func gitCommitAndPush(newBranch string) error {
	// git add all updated files.
	cmd := exec.Command("git", "add", ".")
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
		return ErrNothingChanged
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
	cmd = exec.Command("git", "checkout", "master")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to git checkout master before deleting branch %s: %v", newBranch, err)
	}

	// git branch -D to delete branch to free resources.
	cmd = exec.Command("git", "branch", "-D", newBranch)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to git push branch -D %s: %v", newBranch, err)
	}

	return nil
}

func (g *Generator) sumbitPR(branch string) error {
	title := "docs: auto generate pouch cli/api docs via code"
	head := fmt.Sprintf("pouchrobot:%s", branch)
	base := "master"
	body := `Signed-off-by: pouchrobot <pouch-dev@alibaba-inc.com>

**1.Describe what this PR did**
This PR is automatically done by AI-based collaborating robot.
Pouchrobot will auto-generate cli/api document via https://github.com/spf13/cobra/tree/master/doc every day.
	
**2.Does this pull request fix one issue?** 
None

**3.Describe how you did it**
First, execute command "make client" to build new pouch cli;
Second, execute command "./pouch gen-doc" to generate new cli docs. 

For API part, we use a tool swagger2marup to make it.

**4.Describe how to verify it**
None

**5.Special notes for reviews**
The cli/api doc must be automatically generated.`

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

	return fmt.Sprintf("auto-doc-%s", dateStr)
}

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
	"strings"
	"time"

	"github.com/pouchcontainer/pouchrobot/gh"

	"github.com/google/go-github/github"
	"github.com/sirupsen/logrus"
)

// ErrNothingChanged is used when git commit has nothing to commit.
var ErrNothingChanged = fmt.Errorf("nothing to commit")

// Generator is a processor that periodically auto generated cli and api docs for github repo.
type Generator struct {
	client *gh.Client

	// Owner is the organization of open source project.
	owner string

	// Repo is the repository name.
	repo string

	// RootDir specifies repo's rootdir which is to generated docs.
	RootDir string

	// SwaggerPath specifies that which dir is the swagger.yml file in root dir.
	// this is a relative path to root dir.
	SwaggerPath string

	// APIDocPath specifies where to generate the swagger tool.
	// this is a relative path to root dir.
	APIDocPath string
}

// New initializes a brand new doc generator
func New(client *gh.Client, owner, repo, rootdir, swaggerPath, apiDocPath string) *Generator {
	return &Generator{
		client:      client,
		owner:       owner,
		repo:        repo,
		RootDir:     rootdir,
		SwaggerPath: swaggerPath,
		APIDocPath:  apiDocPath,
	}
}

// Run starts periodical work
// currently generator generates doc every day.
func (g *Generator) Run() error {
	logrus.Infof("start to run doc generator")
	for {
		// Break the loop if the time passes one clock
		// Since time zone container maybe has a delta 8 hours from Beijing Time,
		// It is about 9 o'clock in Beijing Time.
		hour, _, _ := time.Now().Clock()
		logrus.Infof("DocGenerator: now it is %d", hour)
		if hour == 1 {
			break
		}
		time.Sleep(30 * time.Minute)
	}

	for {
		go g.generateDoc()
		// generate cli and api docs every day.
		time.Sleep(24 * time.Hour)
	}
}

// generateDoc starts to generate all docs.
func (g *Generator) generateDoc() error {
	newBranchName := generatenewBranchNameName()
	logrus.Infof("generate a new branch name %s", newBranchName)

	// do prepare thing before cli and api doc generation.
	if err := prepareGitEnv(newBranchName); err != nil {
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

	// auto generate file CONTRIBUTORS on local filesystem.
	if err := g.generateContributors(); err != nil {
		logrus.Errorf("failed to generate CONTRIBUTORS: %v", err)
	}

	// commit and push branch
	if err := gitCommitAndPush(newBranchName); err != nil {
		if err == ErrNothingChanged {
			// if nothing changed, no need to submit pull request.
			return nil
		}
		return err
	}

	// start to submit pull request
	if err := g.sumbitPR(newBranchName); err != nil {
		return err
	}
	return nil
}

func prepareGitEnv(newBranchName string) error {
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

	// create a new branch named by input newBranchName
	// the following doc generation are all on this new branch
	cmd = exec.Command("git", "checkout", "-b", newBranchName)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to git checkout -b %s: %v", newBranchName, err)
	}
	return nil
}

func gitCommitAndPush(newBranchName string) error {
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
	cmd = exec.Command("git", "push", "-f", "origin", newBranchName)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to git push -f origin %s: %v", newBranchName, err)
	}

	// git branch -D to delete branch to free resources.
	cmd = exec.Command("git", "checkout", "master")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to git checkout master before deleting branch %s: %v", newBranchName, err)
	}

	// git branch -D to delete branch to free resources.
	cmd = exec.Command("git", "branch", "-D", newBranchName)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to git push branch -D %s: %v", newBranchName, err)
	}

	return nil
}

func (g *Generator) sumbitPR(branch string) error {
	title := fmt.Sprintf("docs: auto generate %s cli/api docs via code", g.repo)
	head := fmt.Sprintf("pouchrobot:%s", branch)
	base := "master"
	body := `Signed-off-by: pouchrobot <pouch-dev@alibaba-inc.com>

**1.Describe what this PR did**
This PR is automatically done by AI-based collaborating [robot](https://github.com/pouchcontainer/pouchrobot).
Pouchrobot will auto-generate cli/api document via https://github.com/spf13/cobra/tree/master/doc every day.

**2.Does this pull request fix one issue?**
None

**3.Describe how you did it**
First, execute command "make client" to build cli;
Second, execute command "./pouch gen-doc" to generate new cli docs.

For API part, we use a tool swagger2markup to make it.

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

func generatenewBranchNameName() string {
	timeStr := time.Now().String()
	dateStrSlice := strings.SplitN(timeStr, " ", 2)
	dateStr := dateStrSlice[0]

	return fmt.Sprintf("auto-doc-%s", dateStr)
}

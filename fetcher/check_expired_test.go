package fetcher

import (
	"testing"
	"github.com/google/go-github/github"
	"github.com/stretchr/testify/assert"

	"github.com/pouchcontainer/pouchrobot/config"
)

func TestCloseIssue(t *testing.T) {
	type args struct {
		issue *github.Issue
		repo string
		owner string
		token string
	}

	i1 := 8
	i2 := 9
	var issue1 = github.Issue{ID: nil}
	var issue2 = github.Issue{ID: &i1}
	var issue3 = github.Issue{ID: &i2}


	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{name: "test1", args: args{issue: &issue1, repo: "pouch", owner: "test", token: "root"}, wantErr:nil},
		{name: "test2", args: args{issue: &issue2, repo: "pouch", owner: "test", token: "root"}, wantErr:nil},
		{name: "test3", args: args{issue: &issue3, repo: "pouch", owner: "test", token: "root"}, wantErr:nil},

		//{name: "test2", args: args{period: 900},  wantErr: fmt.Errorf("CPU cfs period  %d cannot be less than 1ms (i.e. 1000) or larger than 1s (i.e. 1000000)", 900)},
		//{name: "test3", args: args{period: 1100000},wantErr: fmt.Errorf("CPU cfs period  %d cannot be less than 1ms (i.e. 1000) or larger than 1s (i.e. 1000000)", 1100000)},
		//{name: "test4", args: args{period: 100000}, wantErr: nil},

	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CloseIssue(tt.args.issue, tt.args.repo, tt.args.owner, tt.args.token)
			assert.Equal(t, tt.wantErr, err)

		})
	}
}


func TestCheckExpireIssue(t *testing.T) {
	var fetcher = Fetcher{}
	type args struct {
		config config.Config
	}

	var config1 = config.Config{
		Owner:"test",
		Repo:"pouch",
		AccessToken:"root",
	}

	var config2 = config.Config{
		Owner:"test",
		Repo:"pouch",
		AccessToken:"root",
	}

	var config3 = config.Config{
		Owner:"test",
		Repo:"pouch",
		AccessToken:"root",
	}



	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{name: "test1", args: args{config: config1}, wantErr:nil},
		{name: "test1", args: args{config: config2}, wantErr:nil},
		{name: "test1", args: args{config: config3}, wantErr:nil},

		//{name: "test2", args: args{period: 900},  wantErr: fmt.Errorf("CPU cfs period  %d cannot be less than 1ms (i.e. 1000) or larger than 1s (i.e. 1000000)", 900)},
		//{name: "test3", args: args{period: 1100000},wantErr: fmt.Errorf("CPU cfs period  %d cannot be less than 1ms (i.e. 1000) or larger than 1s (i.e. 1000000)", 1100000)},
		//{name: "test4", args: args{period: 100000}, wantErr: nil},

	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := fetcher.CheckExpireIssue(tt.args.config)
			assert.Equal(t, tt.wantErr, err)

		})
	}
}

package main

import (
	"github.com/pouchcontainer/pouchrobot/server"
	"github.com/pouchcontainer/pouchrobot/server/config"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func main() {
	var cfg config.Config
	var cmdServe = &cobra.Command{
		Use:  "",
		Args: cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			s := server.NewServer(cfg)
			logrus.Fatal(s.Run())
		},
	}

	flagSet := cmdServe.Flags()
	flagSet.StringVarP(&cfg.Owner, "owner", "o", "", "github ID to which connect in GitHub")
	flagSet.StringVarP(&cfg.Repo, "repo", "r", "", "github repo to which connect in GitHub")
	flagSet.StringVarP(&cfg.HTTPListen, "listen", "l", "", "where does automan listened on")
	flagSet.StringVarP(&cfg.AccessToken, "token", "t", "", "access token to have some control on resources")

	cmdServe.Execute()
}

package main

import (
	"github.com/allencloud/automan/server"
	"github.com/allencloud/automan/server/config"

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
	flagSet.StringVarP(&cfg.HTTPListen, "listen", "l", "127.0.0.1:6789", "where does automan listened on")
	flagSet.StringVarP(&cfg.AccessToken, "token", "t", "", "access token to have some control on resources")

	cmdServe.Execute()
}

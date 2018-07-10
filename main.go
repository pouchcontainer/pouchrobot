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

package main

import (
	"github.com/spf13/cobra"
	"github.com/sirupsen/logrus"
	"github.com/pouchcontainer/pouchrobot/config"
)

func main() {
	var cfg config.Config
	var cmdServe = &cobra.Command{
		Use:  "",
		Args: cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			s := NewServer(cfg)
			logrus.Fatal(s.Run())
		},
	}

	flagSet := cmdServe.Flags()
	flagSet.StringVarP(&cfg.Owner, "owner", "o", "paul-yml", "github ID to which connect in GitHub")
	flagSet.StringVarP(&cfg.Repo, "repo", "r", "testrobot", "github repo to which connect in GitHub")
	flagSet.StringVarP(&cfg.HTTPListen, "listen", "l", "", "where does automan listened on")
	flagSet.StringVarP(&cfg.AccessToken, "token", "t", "0a3417b3958a91349625ba0b5d4f4d224ff34d86", "access token to have some control on resources")
	flagSet.StringVarP(&cfg.Timeunit, "timeunit", "u","s", "m for minute, d for day, s for second")
	flagSet.IntVarP(&cfg.Time, "time", "n",1, "use with u, -u d -n 5 means closed issues with out update in 5 days")
	cmdServe.Execute()
}

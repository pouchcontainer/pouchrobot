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
	"github.com/pouchcontainer/pouchrobot/config"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func main() {
	var cfg config.Config
	var cmdServe = &cobra.Command{
		Use:  "An AI-based collaboration robot applied to open source project on GitHub",
		Args: cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			s, err := NewServer(cfg)
			if err != nil {
				logrus.Fatal(err)
			}
			logrus.Fatal(s.Run())
		},
	}
	cmdServe.SilenceUsage = true
	cmdServe.SilenceErrors = true

	flagSet := cmdServe.Flags()
	flagSet.StringVarP(&cfg.Owner, "owner", "o", "", "GitHub username to which this robot connects")
	flagSet.StringVarP(&cfg.Repo, "repo", "r", "", "GitHub code repository to which this robot connects")
	flagSet.StringVarP(&cfg.HTTPListen, "listen", "l", "", "where does robot listen on")
	flagSet.StringVarP(&cfg.AccessToken, "token", "t", "", "access token which identifies robot username in GitHub having write access on repo")
	flagSet.IntVarP(&cfg.CommitsGap, "commits-gap", "c", 20, "commits gap between pull request and master branch; if the fact is beyond this number, request to rebase")

	// for weekly reporter
	flagSet.StringVar(&cfg.ReportDay, "report-day", "Friday", "weekly report generation day of a week")
	flagSet.IntVar(&cfg.ReportHour, "report-hour", 7, "weekly report generation hour on report-day")

	// for doc generator
	flagSet.StringVar(&cfg.RootDir, "root-dir", "", "specifies repo's root directory which is to generated docs")
	flagSet.StringVar(&cfg.SwaggerPath, "swagger-path", "", "specifies where the swagger.yml file locates")
	flagSet.StringVar(&cfg.APIDocPath, "api-doc-path", "", "specifies where to generate the doc file corresponding to swagger.yml")
	flagSet.IntVar(&cfg.GenerationHour, "doc-generation-hour", 1, "specifies doc generation hour of every day.")

	// baidu translator
	flagSet.StringVar(&cfg.BaiduTranslatorAppID, "baidu-trans-appid", "", "specifies the appid for baidu translator")
	flagSet.StringVar(&cfg.BaiduTranslatorKey, "baidu-trans-key", "", "specifies the key for baidu translator")

	cmdServe.Execute()
}

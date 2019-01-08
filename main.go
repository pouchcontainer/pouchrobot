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
	"encoding/json"
	"io/ioutil"

	"github.com/pouchcontainer/pouchrobot/config"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func main() {
	var cmdCfg config.CmdConfig
	cmdExecutor := &cobra.Command{
		Use:  "An AI-based collaboration robot applied to open source project on GitHub",
		Args: cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			configContent, err := ioutil.ReadFile(cmdCfg.ConfigFilePath)
			if err != nil {
				logrus.Fatal(err)
			}

			var cfg config.Config
			if err := json.Unmarshal(configContent, &cfg); err != nil {
				logrus.Fatal(err)
			}

			s, err := NewServer(cfg)
			if err != nil {
				logrus.Fatal(err)
			}

			logrus.Fatal(s.Run())
		},
	}
	cmdExecutor.SilenceUsage = true
	cmdExecutor.SilenceErrors = true

	flagSet := cmdExecutor.Flags()
	flagSet.StringVarP(&cmdCfg.ConfigFilePath, "config", "c", "config.json", "Config file path for robot")

	cmdExecutor.Execute()
}

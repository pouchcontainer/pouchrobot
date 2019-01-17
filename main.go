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
	"os"

	"github.com/pouchcontainer/pouchrobot/config"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var rootCmd *cobra.Command
var cmdCfg config.CmdConfig

func main() {
	rootCmd = &cobra.Command{
		Use:               "pouchrobot",
		Short:             "An AI-based collaboration robot applied to open source project on GitHub",
		Args:              cobra.NoArgs,
		SilenceUsage:      true,
		SilenceErrors:     true,
		DisableAutoGenTag: true, // disable displaying auto generation tag in cli docs
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDaemon(cmd)
		},
	}

	flagSet := rootCmd.Flags()
	flagSet.StringVarP(&cmdCfg.ConfigFilePath, "config", "c", "config.json", "Config file path for robot")
	flagSet.BoolVarP(&cmdCfg.Debug, "debug", "D", false, "Switch daemon log level to DEBUG mode")

	if err := rootCmd.Execute(); err != nil {
		logrus.Error(err)
		os.Exit(1)
	}
}

func runDaemon(cmd *cobra.Command) error {
	if cmdCfg.Debug {
		logrus.Infof("start daemon at debug level")
		logrus.SetLevel(logrus.DebugLevel)
	}

	configContent, err := ioutil.ReadFile(cmdCfg.ConfigFilePath)
	if err != nil {
		return err
	}

	var cfg config.Config
	if err := json.Unmarshal(configContent, &cfg); err != nil {
		return err
	}

	s, err := NewServer(cfg)
	if err != nil {
		return err
	}

	return s.Run()
}

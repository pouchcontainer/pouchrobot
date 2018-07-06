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
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/pouchcontainer/pouchrobot/ci"
	"github.com/pouchcontainer/pouchrobot/config"
	"github.com/pouchcontainer/pouchrobot/docgenerator"
	"github.com/pouchcontainer/pouchrobot/fetcher"
	"github.com/pouchcontainer/pouchrobot/gh"
	"github.com/pouchcontainer/pouchrobot/processor"
	"github.com/pouchcontainer/pouchrobot/reporter"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// DefaultAddress is the default address daemon will listen to.
var DefaultAddress = ":6789"

// Server refers to a daemon server interating with github repos.
type Server struct {
	// listenAddress is the address which is used to accepting requests.
	listenAddress string
	// processor processes webhook event from GitHub.
	processor *processor.Processor
	// fetcher does periodical work to check repo's status on GitHub.
	fetcher *fetcher.Fetcher
	// ciNotifier handles ci system webhook.
	ciNotifier *ci.Notifier
	// reporter reports weekly update of repository.
	reporter *reporter.Reporter
	// docGenerator auto generates docs for repo.
	docGenerator *docgenerator.Generator
}

// NewServer constructs a brand new automan server
func NewServer(config config.Config) *Server {
	ghClient := gh.NewClient(config.Owner, config.Repo, config.AccessToken)
	return &Server{
		listenAddress: config.HTTPListen,
		processor:     processor.New(ghClient),
		fetcher:       fetcher.New(ghClient),
		ciNotifier:    ci.New(ghClient),
		reporter:      reporter.New(ghClient),
		docGenerator:  docgenerator.New(ghClient),
	}
}

// Run runs the server.
func (s *Server) Run() error {
	// start fetcher, reporter and doc generator in goroutines
	go s.fetcher.Run()
	go s.reporter.Run()
	go s.docGenerator.Run()

	// start webserver
	listenAddress := s.listenAddress
	if listenAddress == "" {
		listenAddress = DefaultAddress
	}

	r := mux.NewRouter()

	// register ping api
	r.HandleFunc("/_ping", pingHandler).Methods("GET")

	// github webhook API
	r.HandleFunc("/events", s.gitHubEventHandler).Methods("POST")

	// travisCI webhook API
	r.HandleFunc("/ci_notifications", s.ciNotificationHandler).Methods("POST")

	return http.ListenAndServe(listenAddress, r)
}

// pingHandler handles ping request to return health of server.
func pingHandler(w http.ResponseWriter, r *http.Request) {
	logrus.Debug("/_ping request received")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte{'O', 'K'})
	return
}

// gitHubEventHandler handles webhook events from github.
func (s *Server) gitHubEventHandler(w http.ResponseWriter, r *http.Request) {
	logrus.Debug("/events request received")
	eventType := r.Header.Get("X-Github-Event")

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	r.Body.Close()

	if err := s.processor.HandleEvent(eventType, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}

// ciNotificationHandler handles webhook events from CI system.
func (s *Server) ciNotificationHandler(w http.ResponseWriter, r *http.Request) {
	logrus.Info("/ci_notifications events reveived")
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rawStr := r.PostForm.Get("payload")

	logrus.Debugf("r.PostForm[payload]: %v", rawStr)

	jsonStr := strings.Replace(rawStr, `\"`, `"`, -1)
	if err := s.ciNotifier.Process(jsonStr); err != nil {
		logrus.Errorf("failed to process ci notification: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}

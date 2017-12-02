package server

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/allencloud/automan/server/ci"
	"github.com/allencloud/automan/server/config"
	"github.com/allencloud/automan/server/fetcher"
	"github.com/allencloud/automan/server/gh"
	"github.com/allencloud/automan/server/processor"
	"github.com/allencloud/automan/server/reporter"

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
	}
}

// Run runs the server.
func (s *Server) Run() error {
	// start fetcher and reporter in goroutines
	go s.fetcher.Run()
	go s.reporter.Run()

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

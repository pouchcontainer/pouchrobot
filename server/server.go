package server

import (
	"io/ioutil"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/allencloud/automan/server/config"
	"github.com/allencloud/automan/server/gh"
	"github.com/allencloud/automan/server/processor"
	"github.com/gorilla/mux"
)

// DefaultAddress is the default address daemon will listen to.
var DefaultAddress = ":6789"

// Server refers to a
type Server struct {
	config    config.Config
	processor *processor.Processor
}

// NewServer constructs a brand new automan server
func NewServer(config config.Config) *Server {
	ghClient := gh.NewClient(config.Owner, config.Repo, config.AccessToken)
	return &Server{
		processor: processor.NewProcessor(ghClient),
	}
}

// Run runs the server.
func (s *Server) Run() error {
	listenAddress := s.config.HTTPListen
	if listenAddress == "" {
		listenAddress = DefaultAddress
	}

	r := mux.NewRouter()

	// register ping api
	r.HandleFunc("/_ping", pingHandler).Methods("GET")
	r.HandleFunc("/events", s.eventHandler).Methods("POST")
	return http.ListenAndServe(listenAddress, r)
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	logrus.Info("/_ping request received")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte{'O', 'K'})
	return
}

func (s *Server) eventHandler(w http.ResponseWriter, r *http.Request) {
	logrus.Info("/events request received")
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

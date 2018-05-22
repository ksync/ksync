package syncthing

import (
	"fmt"
	"time"

	"github.com/go-resty/resty"
	log "github.com/sirupsen/logrus"
	"github.com/syncthing/syncthing/lib/config"
	"github.com/syncthing/syncthing/lib/protocol"

	"github.com/vapor-ware/ksync/pkg/debug"
)

var connectRetries = 10


// Server represents a syncthing REST server. It is used to fetch and modify
// configuration as well as restart the syncthing process.
type Server struct {
	Config *config.Configuration `structs:"-"`
	ID     protocol.DeviceID
	URL    string

	client *resty.Client
	stop   chan bool
}

// NewServer constructs a Server with the provided host and apikey.
func NewServer(host string, apikey string) (*Server, error) {
	server := &Server{
		URL:    fmt.Sprintf("http://%s/rest/", host),
		client: resty.New(),
		stop:   make(chan bool),
	}

	server.client.
		SetRetryCount(connectRetries).
		SetRetryWaitTime(1*time.Second).
		SetHostURL(server.URL).
		SetHeader("X-API-KEY", apikey).
		// The resty retry code logs, this makes sure that we log it with our
		// levels and users don't see scary "connection refused" logs by default.
		SetLogger(log.WithFields(log.Fields{}).WriterLevel(log.DebugLevel))

	// TODO: return a friendly error if refresh isn't successful (likely
	// the syncthing server isn't starting up in time) #112
	if err := server.Refresh(); err != nil {
		return nil, err
	}

	return server, nil
}

func (s *Server) String() string {
	return debug.YamlString(s)
}

// Fields returns a set of structured fields for logging.
func (s *Server) Fields() log.Fields {
	return debug.StructFields(s)
}

// Refresh pulls the latest configuration from the configured server and
// updates Server.Config with that value.
func (s *Server) Refresh() error {
	resp, err := s.client.NewRequest().
		SetResult(&config.Configuration{}).
		Get("system/config")
	if err != nil {
		return err
	}

	s.Config = resp.Result().(*config.Configuration)

	id, err := protocol.DeviceIDFromString(resp.Header().Get("X-Syncthing-Id"))
	if err != nil {
		return err
	}
	s.ID = id

	return nil
}

// Update takes the current config, sets the server's config to that and
// restarts the process so that it is live.
func (s *Server) Update() error {
	if _, err := s.client.NewRequest().
		SetBody(s.Config).
		Post("system/config"); err != nil {
		return err
	}

	return s.Restart()
}

// Restart rolls the remote server. Because of how syncthing runs, this just
// happens in the background with only minimal interruption.
func (s *Server) Restart() error {
	if _, err := s.client.NewRequest().Post("system/restart"); err != nil { // nolint: megacheck
		return err
	}

	return nil
}

// Stop is used to clean up everything that this server is doing in the
// background.
func (s *Server) Stop() {
	close(s.stop)
	<-s.stop
	log.WithFields(s.Fields()).Debug("stopping")
}

// IsAlive checks the `ping` endpoint to see if syncthing is up
// TODO: HOTFIX
func (s *Server) IsAlive() (bool, error) {
	time.Sleep(time.Second * 5) // Wait long enough to syncthing has gone down
	if _, err := s.client.NewRequest().Get("system/ping"); err != nil { // nolint: megacheck
		return false, err
	}

	return true, nil
}

package syncthing

import (
	"fmt"
	"time"

	"github.com/go-resty/resty"
	// log "github.com/sirupsen/logrus"
	"github.com/syncthing/syncthing/lib/config"
	"github.com/syncthing/syncthing/lib/protocol"
)

var connectRetries = 10

type Server struct {
	Config *config.Configuration
	ID     protocol.DeviceID

	client *resty.Client
}

func NewServer(host string, apikey string) (*Server, error) {
	server := &Server{}

	server.client = resty.New()
	server.client.
		SetRetryCount(connectRetries).
		SetRetryWaitTime(1*time.Second).
		SetHostURL(fmt.Sprintf("http://%s/rest/", host)).
		SetHeader("X-API-KEY", apikey)

	// TODO: return a friendly error if refresh isn't successful (likely
	// the syncthing server isn't starting up in time) #112
	if err := server.Refresh(); err != nil {
		return nil, err
	}

	return server, nil
}

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

func (s *Server) Update() error {
	if _, err := s.client.NewRequest().
		SetBody(s.Config).
		Post("system/config"); err != nil {
		return err
	}

	return s.Restart()
}

func (s *Server) Restart() error {
	if _, err := s.client.NewRequest().
		Post("system/restart"); err != nil {
		return err
	}

	return nil
}

package syncthing

import (
	"fmt"
	"time"

	"github.com/go-resty/resty"
	log "github.com/sirupsen/logrus"
	"github.com/syncthing/syncthing/lib/config"
	"github.com/syncthing/syncthing/lib/fs"
	"github.com/syncthing/syncthing/lib/protocol"
)

var connectRetries = 5

type Server struct {
	Config *config.Configuration
	ID     protocol.DeviceID

	apikey string
	host   string
}

func NewServer(host string, apikey string) (*Server, error) {
	server := &Server{
		apikey: apikey,
		host:   host,
	}

	if err := server.Refresh(); err != nil {
		return nil, err
	}

	return server, nil
}

func (s *Server) url(path string) string {
	return fmt.Sprintf("http://%s/rest/system/%s", s.host, path)
}

func (s *Server) request() *resty.Request {
	resty.DefaultClient.
		SetRetryCount(connectRetries).
		SetRetryWaitTime(1 * time.Second)
	return resty.R().
		SetHeader("X-API-KEY", s.apikey)
}

func (s *Server) Refresh() error {
	body := &config.Configuration{}

	resp, err := s.request().SetResult(body).Get(s.url("config"))
	if err != nil {
		return err
	}

	s.Config = body

	id, err := protocol.DeviceIDFromString(resp.Header().Get("X-Syncthing-Id"))
	if err != nil {
		return err
	}
	s.ID = id

	return nil
}

func (s *Server) Update() error {
	if _, err := s.request().SetBody(s.Config).Post(s.url("config")); err != nil {
		return err
	}

	return s.Restart()
}

func (s *Server) Restart() error {
	if _, err := s.request().Post(s.url("restart")); err != nil {
		return err
	}

	return nil
}

func Test() {
	q, err := NewServer("127.0.0.1:8384", "ksync")
	if err != nil {
		log.Fatal(err)
	}

	q.Test()
}

func (s *Server) Test() {
	id, err := protocol.DeviceIDFromString(
		"YMMB5TF-LOU47MF-OFLYMHX-PWV6EOI-3EZQFX5-HA2BDL5-BGT2BNQ-ZLHUNQW")
	if err != nil {
		log.Fatal(err)
	}

	device := config.NewDeviceConfiguration(id, "haze.local")
	// device.Addresses = ["127.0.0.1:40000"]

	if err := s.SetDevice(&device); err != nil {
		log.Fatal(err)
	}

	folder := config.NewFolderConfiguration(
		id, "aac3u-twzst", "work-tmp", fs.FilesystemTypeBasic, "/tmp/work")
	folder.FSWatcherEnabled = true
	folder.FSWatcherDelayS = 1

	if err := s.SetFolder(&folder); err != nil {
		log.Fatal(err)
	}

	log.Print(s.Update())
}

package syncthing

import (
	"fmt"

	"github.com/parnurzeal/gorequest"
	log "github.com/sirupsen/logrus"
	"github.com/syncthing/syncthing/lib/config"
	"github.com/syncthing/syncthing/lib/fs"
	"github.com/syncthing/syncthing/lib/protocol"
)

type Server struct {
	Config *config.Configuration

	apikey string
	host   string
}

func NewServer(host string, apikey string) (*Server, []error) {
	server := &Server{
		apikey: apikey,
		host:   host,
	}

	if errs := server.Refresh(); errs != nil {
		return nil, errs
	}

	return server, nil
}

func (s *Server) url(path string) string {
	return fmt.Sprintf("http://%s/rest/system/%s", s.host, path)
}

func (s *Server) get(path string) *gorequest.SuperAgent {
	return gorequest.New().Get(s.url(path)).Set("X-API-KEY", s.apikey)
}

func (s *Server) post(path string) *gorequest.SuperAgent {
	return gorequest.New().Post(s.url(path)).Set("X-API-KEY", s.apikey)
}

func (s *Server) Refresh() []error {
	body := &config.Configuration{}

	if _, _, errs := s.get("config").EndStruct(body); errs != nil {
		return errs
	}

	s.Config = body

	return nil
}

func (s *Server) Update() []error {
	if _, _, errs := s.post("config").Send(s.Config).End(); errs != nil {
		return errs
	}

	return s.Restart()
}

func (s *Server) Restart() []error {
	if _, _, errs := s.post("restart").End(); errs != nil {
		return errs
	}

	return nil
}

func Test() {
	q, errs := NewServer("127.0.0.1:8384", "ksync")
	if errs != nil {
		log.Fatal(errs)
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

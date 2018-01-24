package syncthing

import (
	"fmt"

	"github.com/syncthing/syncthing/lib/config"
	"github.com/syncthing/syncthing/lib/protocol"
)

func (s *Server) GetDevice(id protocol.DeviceID) *config.DeviceConfiguration {
	for _, device := range s.Config.Devices {
		if device.DeviceID == id {
			return &device
		}
	}

	return nil
}

func (s *Server) SetDevice(device *config.DeviceConfiguration) error {
	if err := s.RemoveDevice(device.DeviceID); err != nil {
	}

	s.Config.Devices = append(s.Config.Devices, *device)

	return nil
}

func (s *Server) RemoveDevice(id protocol.DeviceID) error {
	for i, device := range s.Config.Devices {
		if device.DeviceID == id {
			s.Config.Devices[i] = s.Config.Devices[len(s.Config.Devices)-1]
			s.Config.Devices = s.Config.Devices[:len(s.Config.Devices)-1]
			return nil
		}
	}

	return fmt.Errorf("device %s not found", id.String())
}

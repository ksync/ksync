package syncthing

import (
	"fmt"

	"github.com/syncthing/syncthing/lib/config"
)

func (s *Server) GetFolder(id string) *config.FolderConfiguration {
	for _, folder := range s.Config.Folders {
		if folder.ID == id {
			return &folder
		}
	}

	return nil
}

func (s *Server) SetFolder(folder *config.FolderConfiguration) error {
	remoteFolder.FSWatcherEnabled = true
	remoteFolder.FSWatcherDelayS = 1

	if err := s.RemoveFolder(folder.ID); err != nil {
	}

	s.Config.Folders = append(s.Config.Folders, *folder)

	return nil
}

func (s *Server) RemoveFolder(id string) error {
	for i, folder := range s.Config.Folders {
		if folder.ID == id {
			s.Config.Folders[i] = s.Config.Folders[len(s.Config.Folders)-1]
			s.Config.Folders = s.Config.Folders[:len(s.Config.Folders)-1]
			return nil
		}
	}

	return fmt.Errorf("folder %s not found", id)
}

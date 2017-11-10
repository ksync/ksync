package docker

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	log "github.com/sirupsen/logrus"
)

// HasImage checks to see if a docker image exists locally.
func HasImage(name string) (bool, error) {
	args := filters.NewArgs()
	args.Add("reference", name)

	images, err := Client.ImageList(
		context.Background(),
		types.ImageListOptions{
			Filters: args,
		},
	)

	log.WithFields(log.Fields{
		"image": name,
		"count": len(images),
	}).Debug("found image")

	if err != nil {
		return false, err
	}

	if len(images) == 0 {
		return false, nil
	}

	return true, nil
}

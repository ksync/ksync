package server

import (
  "fmt"

  "github.com/golang/protobuf/ptypes/empty"
	"golang.org/x/net/context"
	log "github.com/sirupsen/logrus"

	pb "github.com/vapor-ware/ksync/pkg/proto"
)

func (k *ksyncServer) RestartSyncthing(ctx context.Context, _ *empty.Empty) (*pb.Error, error) {

		if !k.Syncthing.IsAlive() {
			return &pb.Error{Msg: "Syncthing does not appear to be running locally"}, fmt.Errorf("%s", "Syncthing does not appear to be running locally")
		}

		log.Debug("restarting local syncthing")

		return nil, k.Syncthing.Restart()
	}

func (k *ksyncServer) IsAlive(ctx context.Context, _ *empty.Empty) bool {
	return k.Syncthing.IsAlive()
}

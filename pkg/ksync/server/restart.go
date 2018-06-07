package server

import (
	"fmt"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"

	pb "github.com/vapor-ware/ksync/pkg/proto"
)

func (k *ksyncServer) Restart(ctx context.Context, _ *empty.Empty) (*pb.Error, error) {
	log.Warn("requested restart")
	time := time.Second * 10
	k.debounce(ctx, &empty.Empty{}, time)
	return nil, nil
}

func (k *ksyncServer) RestartSyncthing(ctx context.Context, _ *empty.Empty) (*pb.Error, error) {

	if !k.Syncthing.IsAlive() {
		return &pb.Error{Msg: "Syncthing does not appear to be running locally"}, fmt.Errorf("%s", "Syncthing does not appear to be running locally")
	}

	log.Debug("restarting local syncthing")

	return nil, k.Syncthing.Restart()
}

func (k *ksyncServer) IsAlive(ctx context.Context, _ *empty.Empty) (*pb.Alive, error) {
	log.Debug(k.Syncthing)
	switch k.Syncthing.IsAlive() {
	case true:
		return &pb.Alive{Alive: true}, nil
	case false:
		return &pb.Alive{Alive: false}, nil
	}
	return &pb.Alive{Alive: false}, fmt.Errorf("Error during liveness check")
}

func (k *ksyncServer) debounce(ctx context.Context, _ *empty.Empty, t time.Duration) {
	log.Warn("checking debounce")
	incoming := make(chan int)

	go func() {
		var r int

		d := time.NewTimer(t)
		d.Stop()

		for {
			select {
			case r = <-incoming:
				d.Reset(t)
				log.Warn("Got %v requests", r)
			case <-d.C:
				pbErr, err := k.RestartSyncthing(ctx, &empty.Empty{})
				// TODO: This should pass errors on an error channel
				if pbErr != nil {
					log.Fatal(pbErr)
				}
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	}()
}

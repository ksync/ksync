package server

import (
	"fmt"
	"net"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	"github.com/grpc-ecosystem/go-grpc-middleware/tags"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/vapor-ware/ksync/pkg/ksync"
	pb "github.com/vapor-ware/ksync/pkg/proto"
)

type ksyncServer struct {
	SpecList *ksync.SpecList
}

func withDuration(duration time.Duration) (key string, value interface{}) {
	return "grpc.time_ns", duration.Nanoseconds()
}

// Listen starts the ksync server locally.
func Listen(list *ksync.SpecList, bind string, port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", bind, port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.WithFields(log.Fields{
		"bind": bind,
		"port": port,
	}).Info("listening")

	server := &ksyncServer{SpecList: list}

	logrusEntry := log.NewEntry(log.StandardLogger())
	logOpts := []grpc_logrus.Option{
		grpc_logrus.WithDurationField(withDuration),
	}

	opts := []grpc.ServerOption{}

	opts = append(
		opts,
		grpc_middleware.WithUnaryServerChain(
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_logrus.UnaryServerInterceptor(logrusEntry, logOpts...)),
		grpc_middleware.WithStreamServerChain(
			grpc_ctxtags.StreamServerInterceptor(),
			grpc_logrus.StreamServerInterceptor(logrusEntry, logOpts...)),
	)

	rpcServer := grpc.NewServer(opts...)
	pb.RegisterKsyncServer(rpcServer, server)

	return rpcServer.Serve(lis)
}

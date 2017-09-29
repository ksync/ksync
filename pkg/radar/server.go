package radar

import (
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	"github.com/grpc-ecosystem/go-grpc-middleware/tags"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	pb "github.com/vapor-ware/ksync/pkg/proto"
)

type radarServer struct{}

func DefaultServerOpts() []grpc.ServerOption {
	return []grpc.ServerOption{}
}

func withDuration(duration time.Duration) (key string, value interface{}) {
	return "grpc.time_ns", duration.Nanoseconds()
}

// TODO: add readiness/liveliness endpoint (can use prometheus?)
// TODO: add grpc_prometheus
// TODO: add grpc_validator
// TODO: add tracing (net/http/pprof, opentracing?), include in debug logging
func NewServer(opts ...grpc.ServerOption) *grpc.Server {
	logrusEntry := log.NewEntry(log.StandardLogger())
	logOpts := []grpc_logrus.Option{
		grpc_logrus.WithDurationField(withDuration),
	}

	opts = append(
		opts,
		grpc_middleware.WithUnaryServerChain(
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_logrus.UnaryServerInterceptor(logrusEntry, logOpts...)),
		grpc_middleware.WithStreamServerChain(
			grpc_ctxtags.StreamServerInterceptor(),
			grpc_logrus.StreamServerInterceptor(logrusEntry, logOpts...)),
	)

	server := grpc.NewServer(append(DefaultServerOpts(), opts...)...)
	pb.RegisterRadarServer(server, new(radarServer))
	return server
}

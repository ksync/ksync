package cluster

import (
	"fmt"
	"time"

	"github.com/cenkalti/backoff"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/vapor-ware/ksync/pkg/debug"
)

var (
	maxReadyRetries = uint64(10)
)

type Connection struct {
	NodeName string

	service *Service
	tunnels []*Tunnel
}

func NewConnection(nodeName string) *Connection {
	return &Connection{
		NodeName: nodeName,
		service:  NewService(),
		tunnels:  []*Tunnel{},
	}
}

func (c *Connection) String() string {
	return debug.YamlString(c)
}

// Fields returns a set of structured fields for logging.
func (c *Connection) Fields() log.Fields {
	return debug.StructFields(c)
}

// TODO: add TLS
// TODO: add grpc_retry?
func (c *Connection) opts() []grpc.DialOption {
	return []grpc.DialOption{
		grpc.WithTimeout(5 * time.Second),
		grpc.WithBlock(),
		grpc.WithInsecure(),
	}
}

func (c *Connection) waitForHealthy() error {
	test := func() error {
		ready, err := c.service.IsHealthy(c.NodeName)
		if err != nil {
			return backoff.Permanent(err)
		}

		if !ready {
			return fmt.Errorf("radar on %s not ready", c.NodeName)
		}

		return nil
	}

	return backoff.Retry(
		test,
		backoff.WithMaxTries(backoff.NewExponentialBackOff(), maxReadyRetries))
}

func (c *Connection) connection(port int32) (int32, error) {
	if err := c.waitForHealthy(); err != nil {
		return 0, err
	}

	podName, err := c.service.PodName(c.NodeName)
	if err != nil {
		return 0, debug.ErrorOut("cannot get pod name", err, c)
	}

	tun := NewTunnel(c.service.Namespace, podName, port)

	if err := tun.Start(); err != nil {
		return 0, debug.ErrorOut("unable to start tunnel", err, c)
	}

	c.tunnels = append(c.tunnels, tun)

	return tun.LocalPort, nil
}

// Radar creates a new gRPC connection to a radar instance running on
// the specified node.
func (c *Connection) Radar() (*grpc.ClientConn, error) {
	localPort, err := c.connection(c.service.RadarPort)
	if err != nil {
		return nil, debug.ErrorLocation(err)
	}

	return grpc.Dial(fmt.Sprintf("127.0.0.1:%d", localPort), c.opts()...)
}

// SyncthingConnection creates a tunnel to the remote syncthing instance running on
// the specified node.
func (c *Connection) Syncthing() (int32, int32, error) {
	apiPort, err := c.connection(c.service.SyncthingAPI)
	if err != nil {
		return 0, 0, err
	}

	listenerPort, err := c.connection(c.service.SyncthingListener)
	if err != nil {
		return 0, 0, err
	}

	return apiPort, listenerPort, nil
}

func (c *Connection) Stop() error {
	for _, tun := range c.tunnels {
		tun.Close()
	}
	log.WithFields(c.Fields()).Debug("stopped connection")
	return nil
}

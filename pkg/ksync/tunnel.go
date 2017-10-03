package ksync

import (
	"bytes"
	"fmt"
	"net"
	"strconv"

	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/tools/portforward"
	"k8s.io/client-go/tools/remotecommand"
)

type Tunnel struct {
	LocalPort  int
	RemotePort int
	PodName    string
	stopChan   chan struct{}
	readyChan  chan struct{}
	Out        *bytes.Buffer
}

func NewTunnel(nodeName string, remotePort int) (*Tunnel, error) {
	podName, err := radarPodName(nodeName)
	if err != nil {
		return nil, err
	}

	return &Tunnel{
		RemotePort: remotePort,
		PodName:    podName,
		stopChan:   make(chan struct{}, 1),
		readyChan:  make(chan struct{}, 1),
		Out:        new(bytes.Buffer),
	}, nil
}

func (tunnel *Tunnel) Close() {
	close(tunnel.stopChan)
	close(tunnel.readyChan)
}

func (tunnel *Tunnel) Start() error {
	req := KubeClient.CoreV1().RESTClient().Post().
		Resource("pods").
		Namespace(radarNamespace).
		Name(tunnel.PodName).
		SubResource("portforward")

	dialer, err := remotecommand.NewExecutor(KubeCfg, "POST", req.URL())
	if err != nil {
		return err
	}

	local, err := getAvailablePort()
	if err != nil {
		return fmt.Errorf("could not find an available port: %s", err)
	}
	tunnel.LocalPort = local

	log.WithFields(log.Fields{
		"local":  tunnel.LocalPort,
		"remote": tunnel.RemotePort,
		"pod":    tunnel.PodName,
		"url":    req.URL(),
		// TODO: node name?
	}).Debug("starting tunnel")

	pf, err := portforward.New(
		dialer,
		[]string{fmt.Sprintf("%d:%d", tunnel.LocalPort, tunnel.RemotePort)},
		tunnel.stopChan,
		tunnel.readyChan,
		// TODO: there's better places to put this, really anywhere.
		tunnel.Out,
		tunnel.Out)

	if err != nil {
		return err
	}

	errChan := make(chan error)
	go func() {
		errChan <- pf.ForwardPorts()
	}()

	select {
	case err = <-errChan:
		return fmt.Errorf(
			"error forwarding ports (local:%d) (remote:%d) (pod:%s): %v\n%s",
			tunnel.LocalPort,
			tunnel.RemotePort,
			tunnel.PodName,
			err,
			tunnel.Out.String(),
		)
	case <-pf.Ready:
		log.WithFields(log.Fields{
			"local":  tunnel.LocalPort,
			"remote": tunnel.RemotePort,
			"pod":    tunnel.PodName,
			// TODO: node name?
		}).Debug("tunnel running")
		return nil
	}
}

func getAvailablePort() (int, error) {
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		return 0, err
	}
	defer l.Close()

	_, p, err := net.SplitHostPort(l.Addr().String())
	if err != nil {
		return 0, err
	}
	port, err := strconv.Atoi(p)
	if err != nil {
		return 0, err
	}
	return port, err
}

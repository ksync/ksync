package main

import (
	"fmt"
	"net"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/vapor-ware/ksync/pkg/cli"
	"github.com/vapor-ware/ksync/pkg/radar"
)

type serveCmd struct {
	cli.BaseCmd
}

// TODO: update docs, add example.
func (s *serveCmd) new() *cobra.Command {
	long := `Start the server.`
	example := `radar serve --port 40321`

	s.Init("radar", &cobra.Command{
		Use:     "serve [flags]",
		Short:   "Start the server.",
		Long:    long,
		Example: example,
		Args:    cobra.ExactArgs(0),
		Run:     s.run,
	})

	flags := s.Cmd.Flags()
	flags.String(
		"bind",
		"127.0.0.1",
		"interface to which the server will bind")

	if err := s.BindFlag("bind"); err != nil {
		log.Fatal(err)
	}

	flags.IntP(
		"port",
		"p",
		40321,
		"port on which the server will listen")

	if err := s.BindFlag("port"); err != nil {
		log.Fatal(err)
	}

	flags.String("cert", "", "Path to private certificate.")
	if err := s.BindFlag("cert"); err != nil {
		log.Fatal(err)
	}

	flags.String("key", "", "Path to private key.")
	if err := s.BindFlag("key"); err != nil {
		log.Fatal(err)
	}

	return s.Cmd
}

func (s *serveCmd) run(cmd *cobra.Command, args []string) {
	lis, err := net.Listen(
		"tcp", fmt.Sprintf(
			"%s:%d", s.Viper.GetString("bind"), s.Viper.GetInt("port")))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var opts []grpc.ServerOption

	if s.Viper.GetString("cert") != "" && s.Viper.GetString("key") != "" {
		creds, err := credentials.NewServerTLSFromFile(
			s.Viper.GetString("cert"),
			s.Viper.GetString("key"))

		if err != nil {
			log.Fatalf("Failed to generate credentials: %v", err)
		}

		opts = []grpc.ServerOption{grpc.Creds(creds)}
	}

	log.WithFields(log.Fields{
		"bind": s.Viper.GetString("bind"),
		"port": s.Viper.GetInt("port"),
	}).Info("listening")

	server := radar.NewServer(opts...)
	if err := server.Serve(lis); err != nil {
		log.Fatal(err)
	}
}

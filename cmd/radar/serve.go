package main

import (
	"fmt"
	"net"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/vapor-ware/ksync/pkg/radar"
)

type serveCmd struct{}

// TODO: update docs, add example.
func (this *serveCmd) new() *cobra.Command {
	long := `Start the server.`
	example := ``

	cmd := &cobra.Command{
		Use:     "serve [flags]",
		Short:   "Start the server.",
		Long:    long,
		Example: example,
		Args:    cobra.ExactArgs(0),
		Run:     this.run,
	}

	flags := cmd.Flags()
	flags.String(
		"bind",
		"127.0.0.1",
		"interface to which the server will bind")

	viper.BindPFlag("bind", flags.Lookup("bind"))
	viper.BindEnv("bind", "RADAR_BIND")

	flags.IntP(
		"port",
		"p",
		40321,
		"port on which the server will listen")

	viper.BindPFlag("port", flags.Lookup("port"))
	viper.BindEnv("port", "RADAR_PORT")

	flags.String("cert", "", "Path to private certificate.")

	viper.BindPFlag("cert", flags.Lookup("cert"))
	viper.BindEnv("cert", "RADAR_CERT")

	flags.String("key", "", "Path to private key.")

	viper.BindPFlag("key", flags.Lookup("key"))
	viper.BindEnv("key", "RADAR_KEY")

	return cmd
}

func (this *serveCmd) run(cmd *cobra.Command, args []string) {
	lis, err := net.Listen(
		"tcp", fmt.Sprintf("%s:%d", viper.GetString("bind"), viper.GetInt("port")))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var opts []grpc.ServerOption

	if viper.GetString("cert") != "" && viper.GetString("key") != "" {
		creds, err := credentials.NewServerTLSFromFile(
			viper.GetString("cert"),
			viper.GetString("key"))

		if err != nil {
			log.Fatalf("Failed to generate credentials: %v", err)
		}

		opts = []grpc.ServerOption{grpc.Creds(creds)}
	}

	log.WithFields(log.Fields{
		"bind": viper.GetString("bind"),
		"port": viper.GetInt("port"),
	}).Info("listening")

	server := radar.NewServer(opts...)
	server.Serve(lis)
}

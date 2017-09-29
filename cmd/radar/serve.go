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

var (
	serveCmd = &cobra.Command{
		Use:   "serve [flags]",
		Short: "Start the server.",
		Long:  serveDesc,
		Run:   run,
	}

	serveDesc = `Start the server.`
)

func run(cmd *cobra.Command, args []string) {
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

func init() {
	RootCmd.AddCommand(serveCmd)

	serveCmd.Flags().String(
		"bind",
		"127.0.0.1",
		"interface to which the server will bind")

	viper.BindPFlag("bind", serveCmd.Flags().Lookup("bind"))
	viper.BindEnv("bind", "RADAR_BIND")

	serveCmd.Flags().IntP(
		"port",
		"p",
		40321,
		"port on which the server will listen")

	viper.BindPFlag("port", serveCmd.Flags().Lookup("port"))
	viper.BindEnv("port", "RADAR_PORT")

	serveCmd.Flags().String("cert", "", "Path to private certificate.")

	viper.BindPFlag("cert", serveCmd.Flags().Lookup("cert"))
	viper.BindEnv("cert", "RADAR_CERT")

	serveCmd.Flags().String("key", "", "Path to private key.")

	viper.BindPFlag("key", serveCmd.Flags().Lookup("key"))
	viper.BindEnv("key", "RADAR_KEY")
}

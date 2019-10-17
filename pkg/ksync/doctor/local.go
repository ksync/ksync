package doctor

import (
	"fmt"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/ksync/ksync/pkg/ksync"
	pb "github.com/ksync/ksync/pkg/proto"
)

// nolint: staticcheck
var (
	errWatchNotRunning     = fmt.Errorf(`It appears that watch isn't running. You can start it with 'ksync watch'`)
	errWatchNotResponding  = fmt.Errorf(`It appears that watch isn't responding`)
	errSyncthingNotRunning = fmt.Errorf(`It appears that watch isn't running correctly. Please restart 'ksync watch' to correct this`)
)

// DoesSyncthingExist verifies that the local binary exists.
func DoesSyncthingExist() error {
	// There is a timing error when using spinners to output things. If a function
	// completes immediately, you end up with duplicate content. This makes sure
	// that it won't complete immediately.
	time.Sleep(1 * time.Millisecond)

	if !ksync.NewSyncthing().HasBinary() {
		return fmt.Errorf("missing binary, run init to download")
	}

	return nil
}

// IsWatchRunning verifies that watch is running and ready to go.
func IsWatchRunning() error {
	// This is connecting locally and it is very unlikely watch is overloaded,
	// set the timeout *super* short to make it easier on the users when they
	// forgot to start watch.
	withTimeout, _ := context.WithTimeout(context.TODO(), 100*time.Millisecond)

	conn, err := grpc.DialContext(
		withTimeout,
		fmt.Sprintf("127.0.0.1:%d", viper.GetInt("port")),
		[]grpc.DialOption{
			grpc.WithBlock(),
			grpc.WithInsecure(),
		}...)

	if err != nil {
		// The assumption is that the only real error here is because watch isn't
		// running
		log.Debug(err)
		return errWatchNotRunning
	}

	client := pb.NewKsyncClient(conn)
	alive, err := client.IsAlive(context.Background(), &empty.Empty{})
	if err != nil {
		log.Debug(err)
		return errWatchNotResponding
	}

	if !alive.Alive {
		return errSyncthingNotRunning
	}

	if err := conn.Close(); err != nil {
		return err
	}

	return nil
}

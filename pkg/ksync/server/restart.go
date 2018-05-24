package server

import (
  "fmt"

  "github.com/golang/protobuf/ptypes/empty"
	"golang.org/x/net/context"

	"github.com/syncthing/syncthing/lib/rc"

	pb "github.com/vapor-ware/ksync/pkg/proto"
)

func (k *ksyncServer) AwaitAlive(ctx context.Context, _ *empty.Empty) (*pb.Error, error) {
	// @grampelberg: I don't think I'm doing this right. I get a nil pointer
	// error that shows the request is empty?
	// 
	// 	panic: runtime error: invalid memory address or nil pointer dereference
	// [signal SIGSEGV: segmentation violation code=0x1 addr=0x0 pc=0x5068122]
	//
	// goroutine 186 [running]:
	// github.com/vapor-ware/ksync/pkg/ksync/server.(*ksyncServer).AwaitAlive(0xc420246dc0, 0x55a5bc0, 0xc420834570, 0xc420631da0, 0xc420246dc0, 0x5160b20, 0x5588920)
	// 	/Users/timfall/Code/go/src/github.com/vapor-ware/ksync/pkg/ksync/server/restart.go:15 +0x32
	// github.com/vapor-ware/ksync/pkg/proto._Ksync_AwaitAlive_Handler.func1(0x55a5bc0, 0xc420834570, 0x5302a80, 0xc420631da0, 0x0, 0x0, 0x0, 0x0)
	// 	/Users/timfall/Code/go/src/github.com/vapor-ware/ksync/pkg/proto/ksync.pb.go:516 +0x86
	// github.com/vapor-ware/ksync/vendor/github.com/grpc-ecosystem/go-grpc-middleware.ChainUnaryServer.func1.1(0x55a5bc0, 0xc420834570, 0x5302a80, 0xc420631da0, 0x1d, 0xbeb9a4a2aee0f326, 0x735dc68c, 0x5f8e960)

	process := rc.NewProcess(k.Syncthing.URL)

	process.AwaitStartup()

	// if _, err := process.Get("/rest/system/status"); err != nil {
	// 	return nil, err
	// }

	return &pb.Error{Msg: ""}, fmt.Errorf("Something happened")
}

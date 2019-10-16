package server

import (
	"github.com/golang/protobuf/ptypes/empty"
	"golang.org/x/net/context"

	pb "github.com/ksync/ksync/pkg/proto"
)

// GetSpecList returns the list of all registered specs.
func (k *ksyncServer) GetSpecList(
	ctx context.Context, _ *empty.Empty) (*pb.SpecList, error) {

	return k.SpecList.Message()
}

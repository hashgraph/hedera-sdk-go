package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/v2/proto/mirror"
	"github.com/pkg/errors"
	"google.golang.org/grpc/keepalive"
	"time"

	"google.golang.org/grpc"
)

type mirrorNode struct {
	channel *mirror.ConsensusServiceClient
	address string
}

func newMirrorNode(address string) *mirrorNode {
	return &mirrorNode{
		address: address,
		channel: nil,
	}
}

func (node *mirrorNode) getChannel() (*mirror.ConsensusServiceClient, error) {
	if node.channel != nil {
		return node.channel, nil
	}

	var kacp = keepalive.ClientParameters{
		Time:                10 * time.Second,
		Timeout:             time.Second,
		PermitWithoutStream: true,
	}

	conn, err := grpc.Dial(node.address, grpc.WithInsecure(), grpc.WithKeepaliveParams(kacp), grpc.WithBlock())
	if err != nil {
		return nil, errors.Wrapf(err, "error connecting to mirror at %s", node.address)
	}

	channel := mirror.NewConsensusServiceClient(conn)
	node.channel = &channel

	return node.channel, nil
}

package hedera

import (
	"crypto/tls"
	"strings"
	"time"

	"github.com/hashgraph/hedera-protobufs-go/mirror"
	"github.com/pkg/errors"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"

	"google.golang.org/grpc"
)

type mirrorNode struct {
	channel *mirror.ConsensusServiceClient
	client  *grpc.ClientConn
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

	var security grpc.DialOption

	if strings.HasSuffix(node.address, ":50212") || strings.HasSuffix(node.address, ":443") {
		security = grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{}))
	} else {
		security = grpc.WithInsecure()
	}

	conn, err := grpc.Dial(node.address, security, grpc.WithKeepaliveParams(kacp), grpc.WithBlock())
	if err != nil {
		return nil, errors.Wrapf(err, "error connecting to mirror at %s", node.address)
	}

	channel := mirror.NewConsensusServiceClient(conn)
	node.channel = &channel
	node.client = conn

	return node.channel, nil
}

func (node *mirrorNode) close() error {
	if node.channel != nil {
		return node.client.Close()
	}

	return nil
}

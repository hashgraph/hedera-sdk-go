package hedera

import (
	"crypto/tls"
	"time"

	"github.com/hashgraph/hedera-sdk-go/v2/proto/mirror"
	"github.com/pkg/errors"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"

	"google.golang.org/grpc"
)

type _MirrorNode struct {
	channel     *mirror.ConsensusServiceClient
	client      *grpc.ClientConn
	managedNode _ManagedNode
}

func _NewMirrorNode(address string) *_MirrorNode {
	wait := 250 * time.Millisecond
	return &_MirrorNode{
		managedNode: _NewManagedNode(address, wait.Milliseconds()),
		channel:     nil,
	}
}

func (node *_MirrorNode) _SetMinBackoff(waitTime int64) {
	node.managedNode._SetMinBackoff(waitTime)
}

func (node *_MirrorNode) SetAddressBook(addressBook *_NodeAddress) {
	node.managedNode._SetAddressBook(addressBook)
}

func (node *_MirrorNode) GetAddressBook() *_NodeAddress {
	return node.managedNode._GetAddressBook()
}

func (node *_MirrorNode) _InUse() {
	node.managedNode._InUse()
}

func (node *_MirrorNode) _IsHealthy() bool {
	return node.managedNode._IsHealthy()
}

func (node *_MirrorNode) _IncreaseDelay() {
	node.managedNode._IncreaseDelay()
}

func (node *_MirrorNode) _DecreaseDelay() {
	node.managedNode._DecreaseDelay()
}

func (node *_MirrorNode) _Wait() {
	node.managedNode._Wait()
}

func (node *_MirrorNode) _GetChannel() (*mirror.ConsensusServiceClient, error) {
	if node.channel != nil {
		return node.channel, nil
	}

	var kacp = keepalive.ClientParameters{
		Time:                10 * time.Second,
		Timeout:             time.Second,
		PermitWithoutStream: true,
	}

	var security grpc.DialOption

	if node.managedNode.address._IsTransportSecurity() {
		security = grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})) // nolint
	} else {
		security = grpc.WithInsecure()
	}

	conn, err := grpc.Dial(node.managedNode.address._String(), security, grpc.WithKeepaliveParams(kacp), grpc.WithBlock())
	if err != nil {
		return nil, errors.Wrapf(err, "error connecting to mirror at %s", node.managedNode.address._String())
	}

	channel := mirror.NewConsensusServiceClient(conn)
	node.channel = &channel
	node.client = conn

	return node.channel, nil
}

func (node *_MirrorNode) ToSecure() *_MirrorNode {
	node.managedNode.address = node.managedNode.address._ToSecure()
	return node
}

func (node *_MirrorNode) ToInsecure() *_MirrorNode {
	node.managedNode.address = node.managedNode.address._ToInsecure()
	return node
}

func (node *_MirrorNode) _Close() error {
	if node.channel != nil {
		return node.client.Close()
	}

	return nil
}

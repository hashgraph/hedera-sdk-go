package hedera

import (
	"crypto/tls"
	"time"

	"github.com/hashgraph/hedera-protobufs-go/mirror"
	"github.com/pkg/errors"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"

	"google.golang.org/grpc"
)

type _MirrorNode struct {
	*_ManagedNode
	consensusServiceClient *mirror.ConsensusServiceClient
	networkServiceClient   *mirror.NetworkServiceClient
	client                 *grpc.ClientConn
}

func _NewMirrorNode(address string) _MirrorNode {
	wait := 250 * time.Millisecond
	temp := _NewManagedNode(address, wait)
	return _MirrorNode{
		_ManagedNode:           &temp,
		consensusServiceClient: nil,
	}
}

func (node *_MirrorNode) _SetMinBackoff(waitTime time.Duration) {
	node._ManagedNode._SetMinBackoff(waitTime)
}

func (node *_MirrorNode) _GetMinBackoff() time.Duration {
	return node._ManagedNode._GetMinBackoff()
}

func (node *_MirrorNode) _SetMaxBackoff(waitTime time.Duration) {
	node._ManagedNode._SetMaxBackoff(waitTime)
}

func (node *_MirrorNode) _GetMaxBackoff() time.Duration {
	return node._ManagedNode._GetMaxBackoff()
}

func (node *_MirrorNode) _InUse() {
	node._ManagedNode._InUse()
}

func (node *_MirrorNode) _IsHealthy() bool {
	return node._ManagedNode._IsHealthy()
}

func (node *_MirrorNode) _IncreaseDelay() {
	node._ManagedNode._IncreaseDelay()
}

func (node *_MirrorNode) _DecreaseDelay() {
	node._ManagedNode._DecreaseDelay()
}

func (node *_MirrorNode) _Wait() time.Duration {
	return node._ManagedNode._Wait()
}

func (node *_MirrorNode) _GetUseCount() int64 {
	return node._ManagedNode._GetUseCount()
}

func (node *_MirrorNode) _GetLastUsed() int64 {
	return node._ManagedNode._GetLastUsed()
}

func (node *_MirrorNode) _GetManagedNode() *_ManagedNode {
	return node._ManagedNode
}

func (node *_MirrorNode) _GetAttempts() int64 {
	return node._ManagedNode._GetAttempts()
}

func (node *_MirrorNode) _GetAddress() string {
	return node._ManagedNode._GetAddress()
}

func (node *_MirrorNode) _GetConsensusServiceClient() (*mirror.ConsensusServiceClient, error) {
	if node.consensusServiceClient != nil {
		return node.consensusServiceClient, nil
	} else if node.client != nil {
		channel := mirror.NewConsensusServiceClient(node.client)
		node.consensusServiceClient = &channel
		return node.consensusServiceClient, nil
	}

	var kacp = keepalive.ClientParameters{
		Time:                10 * time.Second,
		Timeout:             time.Second,
		PermitWithoutStream: true,
	}

	var security grpc.DialOption

	if node._ManagedNode.address._IsTransportSecurity() {
		security = grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})) // nolint
	} else {
		security = grpc.WithInsecure() //nolint
	}

	conn, err := grpc.Dial(node._ManagedNode.address._String(), security, grpc.WithKeepaliveParams(kacp), grpc.WithBlock())
	if err != nil {
		return nil, errors.Wrapf(err, "error connecting to mirror at %s", node._ManagedNode.address._String())
	}

	channel := mirror.NewConsensusServiceClient(conn)
	node.consensusServiceClient = &channel
	node.client = conn

	return node.consensusServiceClient, nil
}

func (node *_MirrorNode) _GetNetworkServiceClient() (*mirror.NetworkServiceClient, error) {
	if node.networkServiceClient != nil {
		return node.networkServiceClient, nil
	} else if node.client != nil {
		channel := mirror.NewNetworkServiceClient(node.client)
		node.networkServiceClient = &channel
		return node.networkServiceClient, nil
	}

	var kacp = keepalive.ClientParameters{
		Time:                10 * time.Second,
		Timeout:             time.Second,
		PermitWithoutStream: true,
	}

	var security grpc.DialOption

	if node._ManagedNode.address._IsTransportSecurity() {
		security = grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})) // nolint
	} else {
		security = grpc.WithInsecure() //nolint
	}

	conn, err := grpc.Dial(node._ManagedNode.address._String(), security, grpc.WithKeepaliveParams(kacp), grpc.WithBlock())
	if err != nil {
		return nil, errors.Wrapf(err, "error connecting to mirror at %s", node._ManagedNode.address._String())
	}

	channel := mirror.NewNetworkServiceClient(conn)
	node.networkServiceClient = &channel
	node.client = conn

	return node.networkServiceClient, nil
}

func (node *_MirrorNode) _ToSecure() _IManagedNode {
	managed := _ManagedNode{
		address:            node.address._ToSecure(),
		currentBackoff:     node.currentBackoff,
		lastUsed:           node.lastUsed,
		backoffUntil:       node.backoffUntil,
		useCount:           node.useCount,
		minBackoff:         node.minBackoff,
		badGrpcStatusCount: node.badGrpcStatusCount,
	}

	return &_MirrorNode{
		_ManagedNode:           &managed,
		consensusServiceClient: node.consensusServiceClient,
		client:                 node.client,
	}
}

func (node *_MirrorNode) _ToInsecure() _IManagedNode {
	managed := _ManagedNode{
		address:            node.address._ToInsecure(),
		currentBackoff:     node.currentBackoff,
		lastUsed:           node.lastUsed,
		backoffUntil:       node.backoffUntil,
		useCount:           node.useCount,
		minBackoff:         node.minBackoff,
		badGrpcStatusCount: node.badGrpcStatusCount,
	}

	return &_MirrorNode{
		_ManagedNode:           &managed,
		consensusServiceClient: node.consensusServiceClient,
		client:                 node.client,
	}
}

func (node *_MirrorNode) _Close() error {
	if node.consensusServiceClient != nil {
		return node.client.Close()
	}

	return nil
}

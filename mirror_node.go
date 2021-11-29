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
	channel *mirror.ConsensusServiceClient
	client  *grpc.ClientConn
}

func _NewMirrorNode(address string) _MirrorNode {
	wait := 250 * time.Millisecond
	temp := _NewManagedNode(address, wait.Milliseconds())
	return _MirrorNode{
		_ManagedNode: &temp,
		channel:      nil,
	}
}

func (node *_MirrorNode) _SetMinBackoff(waitTime int64) {
	node._ManagedNode._SetMinBackoff(waitTime)
}

func (node *_MirrorNode) _GetMinBackoff() int64 {
	return node._ManagedNode._GetMinBackoff()
}

func (node *_MirrorNode) _SetMaxBackoff(waitTime int64) {
	node._ManagedNode._SetMaxBackoff(waitTime)
}

func (node *_MirrorNode) _GetMaxBackoff() int64 {
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

	if node._ManagedNode.address._IsTransportSecurity() {
		security = grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})) // nolint
	} else {
		security = grpc.WithInsecure()
	}

	conn, err := grpc.Dial(node._ManagedNode.address._String(), security, grpc.WithKeepaliveParams(kacp), grpc.WithBlock())
	if err != nil {
		return nil, errors.Wrapf(err, "error connecting to mirror at %s", node._ManagedNode.address._String())
	}

	channel := mirror.NewConsensusServiceClient(conn)
	node.channel = &channel
	node.client = conn

	return node.channel, nil
}

func (node *_MirrorNode) _ToSecure() _IManagedNode {
	managed := _ManagedNode{
		address:        node.address._ToSecure(),
		currentBackoff: node.currentBackoff,
		lastUsed:       node.lastUsed,
		backoffUntil:   node.lastUsed,
		useCount:       node.useCount,
		minBackoff:     node.minBackoff,
		attempts:       node.attempts,
	}

	return &_MirrorNode{
		_ManagedNode: &managed,
		channel:      node.channel,
		client:       node.client,
	}
}

func (node *_MirrorNode) _ToInsecure() _IManagedNode {
	managed := _ManagedNode{
		address:        node.address._ToInsecure(),
		currentBackoff: node.currentBackoff,
		lastUsed:       node.lastUsed,
		backoffUntil:   node.lastUsed,
		useCount:       node.useCount,
		minBackoff:     node.minBackoff,
		attempts:       node.attempts,
	}

	return &_MirrorNode{
		_ManagedNode: &managed,
		channel:      node.channel,
		client:       node.client,
	}
}

func (node *_MirrorNode) _Close() error {
	if node.channel != nil {
		return node.client.Close()
	}

	return nil
}

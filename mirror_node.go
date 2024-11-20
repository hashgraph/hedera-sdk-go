package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"crypto/tls"
	"time"

	"google.golang.org/grpc/credentials/insecure"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/mirror"
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

func (node *_MirrorNode) _SetVerifyCertificate(_ bool) {
}

func (node *_MirrorNode) _GetVerifyCertificate() bool {
	return false
}

func _NewMirrorNode(address string) (node *_MirrorNode, err error) {
	node = &_MirrorNode{}
	node._ManagedNode, err = _NewManagedNode(address, 250*time.Millisecond)
	return node, err
}

func (node *_MirrorNode) _GetKey() string {
	return node.address._String()
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

func (node *_MirrorNode) _IncreaseBackoff() {
	node._ManagedNode._IncreaseBackoff()
}

func (node *_MirrorNode) _DecreaseBackoff() {
	node._ManagedNode._DecreaseBackoff()
}

func (node *_MirrorNode) _Wait() time.Duration {
	return node._ManagedNode._Wait()
}

func (node *_MirrorNode) _GetUseCount() int64 {
	return node._ManagedNode._GetUseCount()
}

func (node *_MirrorNode) _GetLastUsed() time.Time {
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

func (node *_MirrorNode) _GetReadmitTime() *time.Time {
	return node._ManagedNode._GetReadmitTime()
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
		security = grpc.WithTransportCredentials(insecure.NewCredentials()) //nolint
	}

	conn, err := grpc.NewClient(node._ManagedNode.address._String(), security, grpc.WithKeepaliveParams(kacp))
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
		Time:                time.Minute,
		Timeout:             20 * time.Second,
		PermitWithoutStream: true,
	}

	var security grpc.DialOption

	if node._ManagedNode.address._IsTransportSecurity() {
		security = grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})) // nolint
	} else {
		security = grpc.WithTransportCredentials(insecure.NewCredentials()) //nolint
	}

	conn, err := grpc.NewClient(node._ManagedNode.address._String(), security, grpc.WithKeepaliveParams(kacp))
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
		readmitTime:        node.readmitTime,
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
		readmitTime:        node.readmitTime,
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

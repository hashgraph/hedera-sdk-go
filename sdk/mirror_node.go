package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"math"
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
	consensusServiceClient mirror.ConsensusServiceClient
	networkServiceClient   mirror.NetworkServiceClient
	client                 *grpc.ClientConn
}

func (node *_MirrorNode) _SetVerifyCertificate(_ bool) {
}

func (node *_MirrorNode) _GetVerifyCertificate() bool {
	return false
}

func (node *_MirrorNode) getScheme() (string, error) {
	if node.address.address == nil {
		return "", errors.New("mirror node address is not set")
	}
	port := node.address.port

	// Standard HTTPS ports
	if port == 443 {
		return mirrorNodeSchemeHTTPS, nil
	}

	// Standard HTTP ports
	if port == 80 {
		return mirrorNodeSchemeHTTP, nil
	}

	// For other ports, assume HTTPS for security
	return mirrorNodeSchemeHTTPS, nil
}

func (node *_MirrorNode) getBaseRestUrl() (string, error) {
	scheme, err := node.getScheme()
	if err != nil {
		return "", err
	}
	host := node.address.address
	port := node.address.port

	if host == nil && port == 0 {
		return "", errors.New("mirror node address is not set")
	}
	hostStr := *host
	if hostStr == localhostName || hostStr == "127.0.0.1" {
		return fmt.Sprintf("http://%s:38081/api/v1", hostStr), nil
	}
	return fmt.Sprintf("%s://%s:%d/api/v1", scheme, hostStr, port), nil
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

func (node *_MirrorNode) _GetConsensusServiceClient() (mirror.ConsensusServiceClient, error) {
	if node.consensusServiceClient != nil {
		return node.consensusServiceClient, nil
	} else if node.client != nil {
		channel := mirror.NewConsensusServiceClient(node.client)
		node.consensusServiceClient = channel
		return node.consensusServiceClient, nil
	}

	var kacp = keepalive.ClientParameters{
		Time:                10 * time.Second,
		Timeout:             time.Second,
		PermitWithoutStream: true,
	}

	var security grpc.DialOption

	if node._ManagedNode.address._IsTransportSecurity() {
		security = grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{MinVersion: tls.VersionTLS12}))
	} else {
		security = grpc.WithTransportCredentials(insecure.NewCredentials())
	}

	conn, err := grpc.NewClient(node._ManagedNode.address._String(), security, grpc.WithKeepaliveParams(kacp))
	if err != nil {
		return nil, errors.Wrapf(err, "error connecting to mirror at %s", node._ManagedNode.address._String())
	}

	channel := mirror.NewConsensusServiceClient(conn)
	node.consensusServiceClient = channel
	node.client = conn

	return node.consensusServiceClient, nil
}

func (node *_MirrorNode) _GetNetworkServiceClient() (mirror.NetworkServiceClient, error) {
	if node.networkServiceClient != nil {
		return node.networkServiceClient, nil
	} else if node.client != nil {
		channel := mirror.NewNetworkServiceClient(node.client)
		node.networkServiceClient = channel
		return node.networkServiceClient, nil
	}

	var kacp = keepalive.ClientParameters{
		Time:                time.Minute,
		Timeout:             20 * time.Second,
		PermitWithoutStream: true,
	}

	var security grpc.DialOption

	if node._ManagedNode.address._IsTransportSecurity() {
		security = grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{MinVersion: tls.VersionTLS12}))
	} else {
		security = grpc.WithTransportCredentials(insecure.NewCredentials())
	}

	conn, err := grpc.NewClient(node._ManagedNode.address._String(), security, grpc.WithKeepaliveParams(kacp))
	if err != nil {
		return nil, errors.Wrapf(err, "error connecting to mirror at %s", node._ManagedNode.address._String())
	}

	channel := mirror.NewNetworkServiceClient(conn)
	node.networkServiceClient = channel
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

// RecvStream is a generic interface for any gRPC client stream that has a Recv() method.
// This allows processProtoMessageStream to work with any stream type, not just specific ones.
type recvStream[T any] interface {
	Recv() (T, error)
}

// processProtoMessageStream is a generic method that works with any gRPC client stream
// that implements the RecvStream interface. This allows it to work with any stream type
// such as mirror.NetworkService_GetNodesClient, mirror.ConsensusService_SubscribeTopicClient, etc.
func processProtoMessageStream[T any](ctx context.Context, stream recvStream[T], attempt uint64, maxAttempts uint64, retryHandler func(err error) bool) <-chan streamResult[T] {
	resultStream := make(chan streamResult[T])

	go func() {
		var zero T
		var err error
		defer close(resultStream)
		for {
			select {
			case <-ctx.Done():
				resultStream <- streamResult[T]{zero, ctx.Err()}
				return
			default:
			}
			if err != nil {
				if err == io.EOF {
					return
				}
				if !retryHandler(err) { // should we retry?
					resultStream <- streamResult[T]{zero, err}
					return
				}
				if attempt < maxAttempts { // do we have attempts left?
					delay := math.Min(250.0*math.Pow(2.0, float64(attempt)), 8000)
					time.Sleep(time.Duration(delay) * time.Millisecond)
					attempt++
				} else { // max attempts reached
					resultStream <- streamResult[T]{zero, errors.Wrap(err, "max attempts reached")}
					return
				}
			}

			resp, e := stream.Recv()
			if e != nil {
				err = e
				continue
			}
			resultStream <- streamResult[T]{resp, nil}
		}
	}()

	return resultStream
}

// streamResult is a generic type for stream results
type streamResult[T any] struct {
	data T
	err  error
}

// waitMirrorStreamBackoff waits before a mirror node stream is reopened and reports false if ctx ends first.
func waitMirrorStreamBackoff(ctx context.Context, attempt uint64) bool {
	delay := math.Min(250.0*math.Pow(2.0, float64(attempt)), 8000)

	select {
	case <-ctx.Done():
		return false
	case <-time.After(time.Duration(delay) * time.Millisecond):
		return true
	}
}

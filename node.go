package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"bytes"
	"crypto/sha512"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/status"
)

// _Node represents a node on the network
type _Node struct {
	*_ManagedNode
	accountID         AccountID
	channel           *_Channel
	addressBook       *NodeAddress
	verifyCertificate bool
	channelMutex      sync.Mutex
}

func _NewNode(accountID AccountID, address string, minBackoff time.Duration) (node *_Node, err error) {
	node = &_Node{
		accountID:         accountID,
		verifyCertificate: true,
	}
	node._ManagedNode, err = _NewManagedNode(address, minBackoff)
	return node, err
}

func (node *_Node) _GetKey() string {
	return node.accountID.String()
}

func (node *_Node) _SetMinBackoff(waitTime time.Duration) {
	node._ManagedNode._SetMinBackoff(waitTime)
}

func (node *_Node) _GetMinBackoff() time.Duration {
	return node._ManagedNode._GetMinBackoff()
}

func (node *_Node) _SetMaxBackoff(waitTime time.Duration) {
	node._ManagedNode._SetMaxBackoff(waitTime)
}

func (node *_Node) _GetMaxBackoff() time.Duration {
	return node._ManagedNode._GetMaxBackoff()
}

func (node *_Node) _InUse() {
	node._ManagedNode._InUse()
}

func (node *_Node) _IsHealthy() bool {
	return node._ManagedNode._IsHealthy()
}

func (node *_Node) _IncreaseBackoff() {
	node._ManagedNode._IncreaseBackoff()
}

func (node *_Node) _DecreaseBackoff() {
	node._ManagedNode._DecreaseBackoff()
}

func (node *_Node) _Wait() time.Duration {
	return node._ManagedNode._Wait()
}

func (node *_Node) _GetUseCount() int64 {
	return node._ManagedNode._GetUseCount()
}

func (node *_Node) _GetLastUsed() time.Time {
	return node._ManagedNode._GetLastUsed()
}

func (node *_Node) _GetManagedNode() *_ManagedNode {
	return node._ManagedNode
}

func (node *_Node) _GetAttempts() int64 {
	return node._ManagedNode._GetAttempts()
}

func (node *_Node) _GetAddress() string {
	return node._ManagedNode._GetAddress()
}

func (node *_Node) _GetReadmitTime() *time.Time {
	return node._ManagedNode._GetReadmitTime()
}

func (node *_Node) _GetChannel(logger Logger) (*_Channel, error) {
	node.channelMutex.Lock()
	defer node.channelMutex.Unlock()

	if node.channel != nil {
		return node.channel, nil
	}

	var kacp = keepalive.ClientParameters{
		Time:                10 * time.Second,
		Timeout:             2 * time.Second,
		PermitWithoutStream: true,
	}

	var conn *grpc.ClientConn
	var err error
	security := grpc.WithInsecure() //nolint
	if !node.verifyCertificate {
		println("skipping certificate check")
	}
	if node._ManagedNode.address._IsTransportSecurity() {
		security = grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
			InsecureSkipVerify: true, // nolint
			VerifyPeerCertificate: func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
				if node.addressBook == nil {
					logger.Warn("skipping certificate check since no cert hash was found")
					return nil
				}

				if !node.verifyCertificate {
					return nil
				}

				for _, cert := range rawCerts {
					var certHash []byte

					block := &pem.Block{
						Type:  "CERTIFICATE",
						Bytes: cert,
					}

					var encodedBuf bytes.Buffer
					_ = pem.Encode(&encodedBuf, block)
					digest := sha512.New384()

					if _, err = digest.Write(encodedBuf.Bytes()); err != nil {
						return err
					}

					certHash = digest.Sum(nil)

					if string(node.addressBook.CertHash) == hex.EncodeToString(certHash) {
						return nil
					}
				}

				return x509.CertificateInvalidError{
					Cert:   nil,
					Reason: x509.Expired,
					Detail: "",
				}
			},
		}))
	}

	conn, err = grpc.NewClient(node._ManagedNode.address._String(), security, grpc.WithKeepaliveParams(kacp))
	if err != nil {
		return nil, status.Error(codes.ResourceExhausted, "dial timeout of 10sec exceeded")
	}

	ch := _NewChannel(conn)
	node.channel = &ch

	return node.channel, nil
}

func (node *_Node) _Close() error {
	node.channelMutex.Lock()
	defer node.channelMutex.Unlock()

	if node.channel != nil {
		err := node.channel.client.Close()
		node.channel = nil
		return err
	}

	return nil
}

func (node *_Node) _ToSecure() _IManagedNode {
	managed := _ManagedNode{
		address:            node.address._ToSecure(),
		currentBackoff:     node.currentBackoff,
		lastUsed:           node.lastUsed,
		readmitTime:        node.readmitTime,
		useCount:           node.useCount,
		minBackoff:         node.minBackoff,
		badGrpcStatusCount: node.badGrpcStatusCount,
	}

	return &_Node{
		_ManagedNode:      &managed,
		accountID:         node.accountID,
		channel:           node.channel,
		addressBook:       node.addressBook,
		verifyCertificate: node.verifyCertificate,
	}
}

func (node *_Node) _ToInsecure() _IManagedNode {
	managed := _ManagedNode{
		address:            node.address._ToInsecure(),
		currentBackoff:     node.currentBackoff,
		lastUsed:           node.lastUsed,
		readmitTime:        node.readmitTime,
		useCount:           node.useCount,
		minBackoff:         node.minBackoff,
		badGrpcStatusCount: node.badGrpcStatusCount,
	}

	return &_Node{
		_ManagedNode:      &managed,
		accountID:         node.accountID,
		channel:           node.channel,
		addressBook:       node.addressBook,
		verifyCertificate: node.verifyCertificate,
	}
}

func (node *_Node) _SetVerifyCertificate(verify bool) {
	node.verifyCertificate = verify
}

func (node *_Node) _GetVerifyCertificate() bool {
	return node.verifyCertificate
}

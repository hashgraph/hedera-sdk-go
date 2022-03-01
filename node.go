package hedera

import (
	"bytes"
	"crypto/sha512"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"time"

	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/status"
)

type _Node struct {
	*_ManagedNode
	accountID         AccountID
	channel           *_Channel
	addressBook       *NodeAddress
	verifyCertificate bool
}

type _Nodes struct {
	nodes []_IManagedNode
}

func _NewNode(accountID AccountID, address string, minBackoff int64) _Node {
	temp := _NewManagedNode(address, minBackoff)
	return _Node{
		accountID:         accountID,
		channel:           nil,
		_ManagedNode:      &temp,
		addressBook:       nil,
		verifyCertificate: true,
	}
}

func (node *_Node) _SetMinBackoff(waitTime int64) {
	node._ManagedNode._SetMinBackoff(waitTime)
}

func (node *_Node) _GetMinBackoff() int64 {
	return node._ManagedNode._GetMinBackoff()
}

func (node *_Node) _SetMaxBackoff(waitTime int64) {
	node._ManagedNode._SetMaxBackoff(waitTime)
}

func (node *_Node) _GetMaxBackoff() int64 {
	return node._ManagedNode._GetMaxBackoff()
}

func (node *_Node) _InUse() {
	node._ManagedNode._InUse()
}

func (node *_Node) _IsHealthy() bool {
	return node._ManagedNode._IsHealthy()
}

func (node *_Node) _IncreaseDelay() {
	node._ManagedNode._IncreaseDelay()
}

func (node *_Node) _DecreaseDelay() {
	node._ManagedNode._DecreaseDelay()
}

func (node *_Node) _Wait() time.Duration {
	return node._ManagedNode._Wait()
}

func (node *_Node) _GetUseCount() int64 {
	return node._ManagedNode._GetUseCount()
}

func (node *_Node) _GetLastUsed() int64 {
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

func (node *_Node) _GetChannel() (*_Channel, error) {
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
					println("skipping certificate check since no cert hash was found")
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

	cont, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	conn, err = grpc.DialContext(cont, node._ManagedNode.address._String(), security, grpc.WithKeepaliveParams(kacp), grpc.WithBlock())
	if err != nil {
		return nil, status.Error(codes.ResourceExhausted, "dial timeout of 10sec exceeded")
	}

	ch := _NewChannel(conn)
	node.channel = &ch

	return node.channel, nil
}

func (node *_Node) _Close() error {
	if node.channel != nil {
		err := node.channel.client.Close()
		node.channel = nil
		return err
	}

	return nil
}

func (node *_Node) _ToSecure() _IManagedNode {
	managed := _ManagedNode{
		address:        node.address._ToSecure(),
		currentBackoff: node.currentBackoff,
		lastUsed:       node.lastUsed,
		backoffUntil:   node.lastUsed,
		useCount:       node.useCount,
		minBackoff:     node.minBackoff,
		attempts:       node.attempts,
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
		address:        node.address._ToInsecure(),
		currentBackoff: node.currentBackoff,
		lastUsed:       node.lastUsed,
		backoffUntil:   node.lastUsed,
		useCount:       node.useCount,
		minBackoff:     node.minBackoff,
		attempts:       node.attempts,
	}

	return &_Node{
		_ManagedNode:      &managed,
		accountID:         node.accountID,
		channel:           node.channel,
		addressBook:       node.addressBook,
		verifyCertificate: node.verifyCertificate,
	}
}

func (nodes _Nodes) Len() int {
	return len(nodes.nodes)
}
func (nodes _Nodes) Swap(i, j int) {
	nodes.nodes[i], nodes.nodes[j] = nodes.nodes[j], nodes.nodes[i]
}

func (nodes _Nodes) Less(i, j int) bool {
	if nodes.nodes[i]._IsHealthy() && nodes.nodes[j]._IsHealthy() { // nolint
		if nodes.nodes[i]._GetUseCount() < nodes.nodes[j]._GetUseCount() { // nolint
			return true
		} else if nodes.nodes[i]._GetUseCount() > nodes.nodes[j]._GetUseCount() {
			return false
		} else {
			return nodes.nodes[i]._GetLastUsed() < nodes.nodes[j]._GetLastUsed()
		}
	} else if nodes.nodes[i]._IsHealthy() && !nodes.nodes[j]._IsHealthy() {
		return true
	} else if !nodes.nodes[i]._IsHealthy() && nodes.nodes[j]._IsHealthy() {
		return false
	} else {
		if nodes.nodes[i]._GetUseCount() < nodes.nodes[j]._GetUseCount() { // nolint
			return true
		} else if nodes.nodes[i]._GetUseCount() > nodes.nodes[j]._GetUseCount() {
			return false
		} else {
			return nodes.nodes[i]._GetLastUsed() < nodes.nodes[j]._GetLastUsed()
		}
	}
}

func (node *_Node) _SetCertificateVerification(verify bool) { //nolint
	node.verifyCertificate = verify
}
func (node *_Node) _GetCertificateVerification() bool { //nolint
	return node.verifyCertificate
}

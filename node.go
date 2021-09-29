package hedera

import (
	"bytes"
	"crypto/sha512"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/status"
)

type _Node struct {
	accountID   AccountID
	channel     *_Channel
	managedNode _ManagedNode
}

type _Nodes struct {
	nodes []*_Node
}

func _NewNode(accountID AccountID, address string, minBackoff int64) _Node {
	return _Node{
		accountID:   accountID,
		channel:     nil,
		managedNode: _NewManagedNode(address, minBackoff),
	}
}

func (node *_Node) _SetMinBackoff(waitTime int64) {
	node.managedNode._SetMinBackoff(waitTime)
}

func (node *_Node) SetAddressBook(addressBook *_NodeAddress) {
	node.managedNode._SetAddressBook(addressBook)
}

func (node *_Node) GetAddressBook() *_NodeAddress {
	return node.managedNode._GetAddressBook()
}

func (node *_Node) _InUse() {
	node.managedNode._InUse()
}

func (node *_Node) _IsHealthy() bool {
	return node.managedNode._IsHealthy()
}

func (node *_Node) _IncreaseDelay() {
	node.managedNode._IncreaseDelay()
}

func (node *_Node) _DecreaseDelay() {
	node.managedNode._DecreaseDelay()
}

func (node *_Node) _Wait() {
	node.managedNode._Wait()
}

func (node *_Node) _GetChannel() (*_Channel, error) {
	if node.channel != nil {
		return node.channel, nil
	}

	var kacp = keepalive.ClientParameters{
		Time:                10 * time.Second,
		Timeout:             time.Second,
		PermitWithoutStream: true,
	}

	var conn *grpc.ClientConn
	var err error
	security := grpc.WithInsecure()
	if node.managedNode.address._IsTransportSecurity() {
		security = grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
			InsecureSkipVerify: true, // nolint
			VerifyPeerCertificate: func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
				if node.managedNode.addressBook == nil {
					println("skipping certificate check since no cert hash was found")
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

					if string(node.managedNode.addressBook.certHash) == hex.EncodeToString(certHash) {
						return nil
					}
				}

				return x509.CertificateInvalidError{}
			},
		}))
	}

	cont, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	conn, err = grpc.DialContext(cont, node.managedNode.address._String(), security, grpc.WithKeepaliveParams(kacp), grpc.WithBlock())
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

func (node *_Node) ToSecure() *_Node {
	node.managedNode.address = node.managedNode.address._ToSecure()
	return node
}

func (node *_Node) ToInsecure() *_Node {
	node.managedNode.address = node.managedNode.address._ToInsecure()
	return node
}

func (nodes _Nodes) Len() int {
	return len(nodes.nodes)
}
func (nodes _Nodes) Swap(i, j int) {
	nodes.nodes[i], nodes.nodes[j] = nodes.nodes[j], nodes.nodes[i]
}

func (nodes _Nodes) Less(i, j int) bool {
	if nodes.nodes[i]._IsHealthy() && nodes.nodes[j]._IsHealthy() { // nolint
		if nodes.nodes[i].managedNode.useCount < nodes.nodes[j].managedNode.useCount { // nolint
			return true
		} else if nodes.nodes[i].managedNode.useCount > nodes.nodes[j].managedNode.useCount {
			return false
		} else {
			return nodes.nodes[i].managedNode.lastUsed < nodes.nodes[j].managedNode.lastUsed
		}
	} else if nodes.nodes[i]._IsHealthy() && !nodes.nodes[j]._IsHealthy() {
		return true
	} else if !nodes.nodes[i]._IsHealthy() && nodes.nodes[j]._IsHealthy() {
		return false
	} else {
		if nodes.nodes[i].managedNode.useCount < nodes.nodes[j].managedNode.useCount { // nolint
			return true
		} else if nodes.nodes[i].managedNode.useCount > nodes.nodes[j].managedNode.useCount {
			return false
		} else {
			return nodes.nodes[i].managedNode.lastUsed < nodes.nodes[j].managedNode.lastUsed
		}
	}
}

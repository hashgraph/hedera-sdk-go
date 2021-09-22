package hedera

import (
	"bytes"
	"crypto/sha512"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"math"
	"strings"
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
	address     string
	delay       int64
	lastUsed    int64
	delayUntil  int64
	useCount    int64
	channel     *_Channel
	waitTime    int64
	attempts    int64
	addressBook *_NodeAddress
}

type _Nodes struct {
	nodes []*_Node
}

func _NewNode(accountID AccountID, address string, waitTime int64) _Node {
	return _Node{
		accountID:   accountID,
		address:     address,
		delay:       250,
		lastUsed:    time.Now().UTC().UnixNano(),
		delayUntil:  time.Now().UTC().UnixNano(),
		useCount:    0,
		channel:     nil,
		waitTime:    waitTime,
		attempts:    0,
		addressBook: nil,
	}
}

func (node *_Node) _SetWaitTime(waitTime int64) {
	if node.delay == node.waitTime {
		node.delay = node.waitTime
	}

	node.waitTime = waitTime
}

func (node *_Node) SetAddressBook(addressBook *_NodeAddress) {
	node.addressBook = addressBook
}

func (node *_Node) GetAddressBook() *_NodeAddress {
	return node.addressBook
}

func (node *_Node) _InUse() {
	node.useCount++
	node.lastUsed = time.Now().UTC().UnixNano()
}

func (node *_Node) _IsHealthy() bool {
	return node.delayUntil <= time.Now().UTC().UnixNano()
}

func (node *_Node) _IncreaseDelay() {
	node.attempts++
	node.delay = int64(math.Min(float64(node.delay)*2, 8000))
	node.delayUntil = (node.delay * 100000) + time.Now().UTC().UnixNano()
}

func (node *_Node) _DecreaseDelay() {
	node.delay = int64(math.Max(float64(node.delay)/2, 250))
}

func (node *_Node) _Wait() {
	delay := node.delayUntil - node.lastUsed
	time.Sleep(time.Duration(delay) * time.Nanosecond)
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
	parts := strings.SplitN(node.address, ":", 2)
	security := grpc.WithInsecure()
	if parts[1] == "443" || parts[1] == "50212" {
		security = grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
			InsecureSkipVerify: true, // nolint
			VerifyPeerCertificate: func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
				if node.addressBook == nil {
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

					if string(node.addressBook.certHash) == hex.EncodeToString(certHash) {
						return nil
					}
				}

				return x509.CertificateInvalidError{}
			},
		}))
	}

	cont, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	conn, err = grpc.DialContext(cont, node.address, security, grpc.WithKeepaliveParams(kacp), grpc.WithBlock())
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

func (nodes _Nodes) Len() int {
	return len(nodes.nodes)
}
func (nodes _Nodes) Swap(i, j int) {
	nodes.nodes[i], nodes.nodes[j] = nodes.nodes[j], nodes.nodes[i]
}

func (nodes _Nodes) Less(i, j int) bool {
	if nodes.nodes[i]._IsHealthy() && nodes.nodes[j]._IsHealthy() { // nolint
		if nodes.nodes[i].useCount < nodes.nodes[j].useCount { // nolint
			return true
		} else if nodes.nodes[i].useCount > nodes.nodes[j].useCount {
			return false
		} else {
			return nodes.nodes[i].lastUsed < nodes.nodes[j].lastUsed
		}
	} else if nodes.nodes[i]._IsHealthy() && !nodes.nodes[j]._IsHealthy() {
		return true
	} else if !nodes.nodes[i]._IsHealthy() && nodes.nodes[j]._IsHealthy() {
		return false
	} else {
		if nodes.nodes[i].useCount < nodes.nodes[j].useCount { // nolint
			return true
		} else if nodes.nodes[i].useCount > nodes.nodes[j].useCount {
			return false
		} else {
			return nodes.nodes[i].lastUsed < nodes.nodes[j].lastUsed
		}
	}
}

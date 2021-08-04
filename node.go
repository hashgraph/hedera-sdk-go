package hedera

import (
	"crypto/tls"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/status"
	"math"
	"strings"
	"time"
)

type node struct {
	accountID  AccountID
	address    string
	delay      int64
	lastUsed   int64
	delayUntil int64
	useCount   int64
	channel    *channel
	waitTime   int64
	attempts   int64
}

type nodes struct {
	nodes []*node
}

func newNode(accountID AccountID, address string, waitTime int64) node {
	return node{
		accountID:  accountID,
		address:    address,
		delay:      250,
		lastUsed:   time.Now().UTC().UnixNano(),
		delayUntil: time.Now().UTC().UnixNano(),
		useCount:   0,
		channel:    nil,
		waitTime:   waitTime,
		attempts:   0,
	}
}

func (node *node) setWaitTime(waitTime int64) {
	if node.delay == node.waitTime {
		node.delay = node.waitTime
	}

	node.waitTime = waitTime
}

func (node *node) inUse() {
	node.useCount += 1
	node.lastUsed = time.Now().UTC().UnixNano()
}

func (node *node) isHealthy() bool {
	return node.delayUntil <= time.Now().UTC().UnixNano()
}

func (node *node) increaseDelay() {
	node.attempts++
	node.delay = int64(math.Min(float64(node.delay)*2, 8000))
	node.delayUntil = (node.delay * 100000) + time.Now().UTC().UnixNano()
}

func (node *node) decreaseDelay() {
	node.delay = int64(math.Max(float64(node.delay)/2, 250))
}

func (node *node) wait() {
	delay := node.delayUntil - node.lastUsed
	time.Sleep(time.Duration(delay) * time.Nanosecond)

}

func (node *node) getChannel() (*channel, error) {
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
		security = grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{InsecureSkipVerify: true}))
	}

	cont, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	conn, err = grpc.DialContext(cont, node.address, security, grpc.WithKeepaliveParams(kacp), grpc.WithBlock())
	if err != nil {
		return nil, status.Error(codes.ResourceExhausted, "dial timeout of 10sec exceeded")
	}

	ch := newChannel(conn)
	node.channel = &ch

	return node.channel, nil
}

func (node *node) close() error {
	if node.channel != nil {
		err := node.channel.client.Close()
		node.channel = nil
		return err
	}

	return nil
}

func (nodes nodes) Len() int {
	return len(nodes.nodes)
}
func (nodes nodes) Swap(i, j int) {
	nodes.nodes[i], nodes.nodes[j] = nodes.nodes[j], nodes.nodes[i]
}

func (nodes nodes) Less(i, j int) bool {
	if nodes.nodes[i].isHealthy() && nodes.nodes[j].isHealthy() {
		if nodes.nodes[i].useCount < nodes.nodes[j].useCount {
			return true
		} else if nodes.nodes[i].useCount > nodes.nodes[j].useCount {
			return false
		} else {
			return nodes.nodes[i].lastUsed < nodes.nodes[j].lastUsed
		}
	} else if nodes.nodes[i].isHealthy() && !nodes.nodes[j].isHealthy() {
		return true
	} else if !nodes.nodes[i].isHealthy() && nodes.nodes[j].isHealthy() {
		return false
	} else {
		if nodes.nodes[i].useCount < nodes.nodes[j].useCount {
			return true
		} else if nodes.nodes[i].useCount > nodes.nodes[j].useCount {
			return false
		} else {
			return nodes.nodes[i].lastUsed < nodes.nodes[j].lastUsed
		}
	}
}

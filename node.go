package hedera

import (
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"math"
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
}

type nodes struct {
	nodes []*node
}

func newNode(accountID AccountID, address string) node {
	return node{
		accountID:  accountID,
		address:    address,
		delay:      250,
		lastUsed:   time.Now().UTC().UnixNano(),
		delayUntil: time.Now().UTC().UnixNano(),
		useCount:   0,
		channel:    nil,
	}
}

func (node *node) inUse() {
	node.useCount += 1
	node.lastUsed = time.Now().UTC().UnixNano()
}

func (node *node) isHealthy() bool {
	return node.delayUntil <= time.Now().UTC().UnixNano()
}

func (node *node) increaseDelay() {
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

	conn, err := grpc.Dial(node.address, grpc.WithInsecure(), grpc.WithKeepaliveParams(kacp), grpc.WithBlock())
	if err != nil {
		return nil, errors.Wrapf(err, "error connecting to node at %s", node.address)
	}

	ch := newChannel(conn)
	node.channel = &ch

	return node.channel, nil
}

func (node *node) close() error {
	return node.channel.client.Close()
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

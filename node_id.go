package hedera

import (
	"math"
	"time"
)

type NodeID struct {
	AccountID AccountID
	Address string
	delay int64
	lastUsed *int64
}

func NewNodeID(accountID AccountID, address string) NodeID {
	return NodeID{
		AccountID: accountID,
		Address:   address,
		delay:     250,
		lastUsed:  nil,
	}
}

func (node *NodeID) isHealthy() bool {
	lastUsed := *node.lastUsed
	if node.lastUsed != nil {
		return lastUsed + node.delay < time.Now().UTC().UnixNano() / 1e6
	}

	return true
}

func (node *NodeID) increaseDelay(){
	*node.lastUsed = time.Now().UTC().UnixNano()/1e6
	node.delay = int64(math.Min(float64(node.delay) * 2, 8000))
}

func (node *NodeID) decreaseDelay(){
	node.delay = int64(math.Max(float64(node.delay) / 2, 250))
}

func (node *NodeID) wait(){
	delay := *node.lastUsed
	if node.lastUsed != nil {
		delay = delay + node.delay - time.Now().UTC().UnixNano() / 1e6
	} else {
		delay = node.delay - time.Now().UTC().UnixNano() / 1e6
	}

	time.Sleep(time.Duration(delay) * time.Millisecond)
}

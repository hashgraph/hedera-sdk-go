package hedera

import (
	"math"
	"time"
)

type _ManagedNode struct {
	accountID      AccountID
	address        *_ManagedNodeAddress
	currentBackoff int64
	lastUsed       int64
	backoffUntil   int64
	useCount       int64
	minBackoff     int64
	attempts       int64
	addressBook    *_NodeAddress
}

func _NewManagedNode(address string, minBackoff int64) _ManagedNode {
	return _ManagedNode{
		address:        _ManagedNodeAddressFromString(address),
		currentBackoff: 250,
		lastUsed:       time.Now().UTC().UnixNano(),
		backoffUntil:   time.Now().UTC().UnixNano(),
		useCount:       0,
		minBackoff:     minBackoff,
		attempts:       0,
		addressBook:    nil,
	}
}

func (node *_ManagedNode) _SetMinBackoff(waitTime int64) {
	if node.currentBackoff == node.minBackoff {
		node.currentBackoff = node.minBackoff
	}

	node.minBackoff = waitTime
}

func (node *_ManagedNode) _SetAddressBook(addressBook *_NodeAddress) {
	node.addressBook = addressBook
}

func (node *_ManagedNode) _GetAddressBook() *_NodeAddress {
	return node.addressBook
}

func (node *_ManagedNode) _InUse() {
	node.useCount++
	node.lastUsed = time.Now().UTC().UnixNano()
}

func (node *_ManagedNode) _IsHealthy() bool {
	return node.backoffUntil <= time.Now().UTC().UnixNano()
}

func (node *_ManagedNode) _IncreaseDelay() {
	node.attempts++
	node.currentBackoff = int64(math.Min(float64(node.currentBackoff)*2, 8000))
	node.backoffUntil = (node.currentBackoff * 100000) + time.Now().UTC().UnixNano()
}

func (node *_ManagedNode) _DecreaseDelay() {
	node.currentBackoff = int64(math.Max(float64(node.currentBackoff)/2, 250))
}

func (node *_ManagedNode) _Wait() {
	delay := node.backoffUntil - node.lastUsed
	time.Sleep(time.Duration(delay) * time.Nanosecond)
}

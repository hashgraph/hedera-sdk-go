package hedera

import (
	"math"
	"time"
)

type _IManagedNode interface {
	_SetMinBackoff(waitTime int64)
	_InUse()
	_IsHealthy() bool
	_IncreaseDelay()
	_DecreaseDelay()
	_Wait()
	_GetUseCount() int64
	_GetLastUsed() int64
	_GetAttempts() int64
	_GetAddress() string
	_ToSecure() _IManagedNode
	_ToInsecure() _IManagedNode
	_GetManagedNode() *_ManagedNode
	_Close() error
}

type _ManagedNode struct {
	address        *_ManagedNodeAddress
	currentBackoff int64
	lastUsed       int64
	backoffUntil   int64
	useCount       int64
	minBackoff     int64
	attempts       int64
	inUse          bool
}

func (node *_ManagedNode) _GetAttempts() int64 {
	return node.attempts
}

func (node *_ManagedNode) _GetAddress() string {
	if node.address != nil {
		return node.address._String()
	}

	return ""
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
		inUse:          false,
	}
}

func (node *_ManagedNode) _SetMinBackoff(waitTime int64) {
	if node.currentBackoff == node.minBackoff {
		node.currentBackoff = node.minBackoff
	}

	node.minBackoff = waitTime
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

func (node *_ManagedNode) _GetUseCount() int64 {
	return node.useCount
}

func (node *_ManagedNode) _GetLastUsed() int64 {
	return node.lastUsed
}
